package routes

import (
	"authentication-system/controllers"
	"authentication-system/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController *controllers.AuthController
}

func NewAuthRouteController(authController *controllers.AuthController) *AuthRouteController {
	return &AuthRouteController{authController: authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.GET("me",middleware.AuthMiddleware(), rc.authController.GetMe)
	router.POST("/register", rc.authController.RegisterUser)
	router.POST("/login", rc.authController.LoginUser)
	router.POST("/logout", rc.authController.LogoutUser)
	router.POST("/refresh", rc.authController.RefreshAccessToken)

}
