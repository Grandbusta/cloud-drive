package main

import (
	"fmt"
	"log"

	"github.com/Grandbusta/cloud-drive/config"
	"github.com/Grandbusta/cloud-drive/controllers"
	"github.com/Grandbusta/cloud-drive/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println(".env loaded")
	}
}

func main() {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	db := config.NewDB()
	db.AutoMigrate()

	userRoutes := r.Group("/user")
	userRoutes.POST("/signup", controllers.CreateUser)

	r.Run(":8080")
}
