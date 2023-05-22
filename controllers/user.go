package controllers

import (
	"net/http"

	"github.com/Grandbusta/cloud-drive/utils"
	"github.com/gin-gonic/gin"
)

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(ctx *gin.Context) {
	var user CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	if isValid := utils.ValidEmail(user.Email); !isValid {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	password, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	user.Password = password
	utils.SuccessWithMessage(ctx, http.StatusOK, "endpoint successful")
}
