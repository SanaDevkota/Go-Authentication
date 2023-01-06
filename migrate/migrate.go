package main

import (
	"authentication-system/config"
	"authentication-system/models"
	"fmt"
	"log"
)

func init() {
	configuration, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	config.ConnectDB(&configuration)

}

func main() {
	config.DB.AutoMigrate(&models.User{})
	fmt.Println("👍 Migration complete")

}
