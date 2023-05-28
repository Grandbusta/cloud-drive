package main

import (
	"fmt"
	"log"

	"github.com/Grandbusta/cloud-drive/config"
	"github.com/Grandbusta/cloud-drive/controllers"
	"github.com/Grandbusta/cloud-drive/middlewares"
	"github.com/Grandbusta/cloud-drive/models"
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
	db.Debug().AutoMigrate(&models.User{}, &models.Resource{}, &models.TreePath{})

	r.MaxMultipartMemory = 8 << 20
	userRoutes := r.Group("/user")
	userRoutes.POST("/signup", controllers.CreateUser)
	userRoutes.POST("/signin", controllers.LoginUser)

	resourceRoutes := r.Group("/resource")
	resourceRoutes.POST("/create-folder", middlewares.TokenAuthMiddleware(), controllers.CreateFolder)
	resourceRoutes.POST("/upload-file", middlewares.TokenAuthMiddleware(), controllers.UploadFile)
	resourceRoutes.POST("/update/:resource_id", middlewares.TokenAuthMiddleware(), controllers.UpdateResource)
	resourceRoutes.GET("/delete/:resource_id", middlewares.TokenAuthMiddleware(), controllers.DeleteResource)
	resourceRoutes.GET("/data/:resource_id", middlewares.TokenAuthMiddleware(), controllers.GetResource)

	r.Run(":8080")
}
