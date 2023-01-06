package main

import (
	"authentication-system/config"
	"authentication-system/controllers"
	"authentication-system/routes"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	// PostController      controllers.PostController
	// PostRouteController routes.PostRouteController
)

func init() {
	//loading config
	envconfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	// connecting to database
	config.ConnectDB(&envconfig)

	AuthController = *controllers.NewAuthController(config.DB)
	AuthRouteController = *routes.NewAuthRouteController(&AuthController)

	UserController = *controllers.NewUserController(config.DB)
	UserRouteController = *routes.NewUserRouteController(UserController)
	// initializing gin server with default middleware
	server = gin.Default()
}

func main() {
	//loading config
	envconfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	router := server.Group("/api")

	router.GET("/", func(ctx *gin.Context) {
		message := "Welcome to the Authentication System"
		ctx.JSON(200, gin.H{
			"success": true,
			"message": message,
		},
		)
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)

	log.Fatal(server.Run(":" + envconfig.ServerPort))

}
