package controllers

import (
	"authentication-system/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User

	query := uc.DB.Where("id = ?", id).First(&user)
	if query.Error != nil {
		ctx.JSON(404, gin.H{"success": false, "message": "User not found"})
		return
	}
	ctx.JSON(200, gin.H{"success": true, "data": user})


}


func (uc *UserController) GetUsers(ctx *gin.Context) {
	var users []models.User
	query := uc.DB.Find(&users)

	if query.Error != nil {
		ctx.JSON(404, gin.H{"success": false, "message": "No users found"})
		return
	}

	ctx.JSON(200, gin.H{"success": true, "data": users})
}


