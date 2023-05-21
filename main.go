package main

import (
	"github.com/Grandbusta/cloud-drive/controllers"
	"github.com/Grandbusta/cloud-drive/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	userRoutes := r.Group("/user")
	userRoutes.POST("/signup", controllers.CreateUser)

	r.Run(":8080")
}
