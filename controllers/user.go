package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Grandbusta/cloud-drive/config"
	"github.com/Grandbusta/cloud-drive/models"
	"github.com/Grandbusta/cloud-drive/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(ctx *gin.Context) {
	var userInput CreateUserInput
	db := config.NewDB()
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	if isValid := utils.ValidEmail(userInput.Email); !isValid {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	password, err := utils.HashPassword(userInput.Password)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	user := models.User{}
	user.Password = password
	user.Email = userInput.Email
	existingUser, err := user.FindUserByEmail(db)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	if existingUser.Email != "" {
		utils.ServerResponse(ctx, http.StatusConflict, "User already exist")
		return
	}

	newUser, err := user.SaveUser(db)
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithMessageeAndData(
		ctx,
		http.StatusCreated,
		"User created successfully",
		newUser,
	)
}
