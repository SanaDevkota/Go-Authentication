package controllers

import (
	"authentication-system/config"
	"authentication-system/models"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"authentication-system/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Register a User
func (ac *AuthController) RegisterUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBind(&payload); err != nil {
		log.Print(err)
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if payload.Password != payload.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.Hashpassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "Password deosnot match"})
		return
	}

	now := time.Now()
	newuser := models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	result := ac.DB.Create(&newuser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"success": false, "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Something bad happened"})
		return
	}

	userResponse := models.UserResponse{
		ID:        uint(newuser.ID),
		Name:      newuser.Name,
		Email:     newuser.Email,
		Role:      newuser.Role,
		CreatedAt: newuser.CreatedAt,
		UpdatedAt: newuser.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "data": gin.H{"user": userResponse}})

}

// Login a User
func (ac *AuthController) LoginUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User

	result := ac.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email or Password"})
		return
	}

	envconfig, _ := config.LoadConfig(".")

	access_token, err := utils.CreateToken(envconfig.AccessTokenExpiresIn, user.ID, envconfig.AccessTokenPrivateKey)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	refresh_token, err := utils.CreateToken(envconfig.RefreshTokenExpiresIn, user.ID, envconfig.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.SetCookie("access_token", access_token, envconfig.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.SetCookie("refresh_token", refresh_token, envconfig.RefreshTokenMaxAge*60, "/", "localhost", false, false)
	ctx.SetCookie("logged_in", "true", envconfig.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"success": true, "access_token": access_token})

}

// Logout a User
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, false)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, false)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// Get the current logged in user
func (ac *AuthController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User) // type assertion with models.user

	userResponse := &models.UserResponse{
		ID:        uint(currentUser.ID),
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Role:      currentUser.Role,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})

}

// Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	refresh_token, err := ctx.Cookie("refresh_token")
	if err != nil {
		message := "You are not logged in"
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": message})
		return
	}
	//loading the environment variables
	envconfig, _ := config.LoadConfig(".")

	// validating the refresh token
	sub, err := utils.ValidateToken(refresh_token, envconfig.RefreshTokenPublicKey)
	if err != nil {

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "the user belonging to this token no logger exists"})

		return
	}

	var user models.User
	// sub means subject and in subject we have the user id so we are getting the user from the database
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "the user belonging to this token no logger exists"})
		return
	}
	access_token, err := utils.CreateToken(envconfig.AccessTokenExpiresIn, user.ID, envconfig.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": true, "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, envconfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", envconfig.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"success": true, "access_token": access_token})

}
