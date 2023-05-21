package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserInput struct {
	Email    string
	Password string
}

func CreateUser(ctx *gin.Context) {
	var user CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "invalid-json",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "sender",
	})
}
