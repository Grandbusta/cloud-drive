package utils

import "github.com/gin-gonic/gin"

func ServerResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}
