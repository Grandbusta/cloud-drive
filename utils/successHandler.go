package utils

import "github.com/gin-gonic/gin"

func SuccessWithMessage(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
	})
}

func SuccessWithData(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

func SuccessWithMessageAndData(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"data":    data,
		"message": message,
	})
}
