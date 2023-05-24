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

func CreateUser(ctx *gin.Context) {
	var userInput models.CreateUserInput
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
	utils.SuccessWithMessageAndData(
		ctx,
		http.StatusCreated,
		"User created successfully",
		newUser.PublicUser(),
	)
}

func LoginUser(ctx *gin.Context) {
	var userInput models.LoginUserInput
	db := config.NewDB()
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	user := models.User{}
	user.Email = userInput.Email
	existingUser, err := user.FindUserByEmail(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ServerResponse(ctx, http.StatusNotFound, "User not found")
		return
	}
	if err := utils.ComparePassword(existingUser.Password, userInput.Password); err != nil {
		utils.ServerResponse(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	token, err := utils.CreateToken(existingUser.ID)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithData(ctx, http.StatusOK, map[string]interface{}{
		"user":  existingUser.PublicUser(),
		"token": token,
	},
	)
}
