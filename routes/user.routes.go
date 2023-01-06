package routes

import (
	"authentication-system/controllers"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) *UserRouteController {
	return &UserRouteController{userController: userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.GET("/:id", uc.userController.GetUser)
	router.GET("/", uc.userController.GetUsers)

}
