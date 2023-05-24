package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Grandbusta/cloud-drive/config"
	"github.com/Grandbusta/cloud-drive/models"
	"github.com/Grandbusta/cloud-drive/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateFolder(ctx *gin.Context) {
	var folderInput models.CreateFolderInput
	db := config.NewDB()
	userID, err := utils.ExtractTokenId(ctx)
	if err != nil || userID == "" {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}

	if err := ctx.ShouldBindJSON(&folderInput); err != nil {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	resource := models.Resource{}
	fmt.Println(folderInput)
	resource.ParentID = folderInput.ParentID
	if folderInput.ParentID != config.ROOT {
		parentResource, err := resource.FindResourceByID(db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ServerResponse(ctx, http.StatusNotFound, "parent_id not found")
			return
		}
		if parentResource.ResourceType != config.RESOURCE_TYPE_FOLDER {
			utils.ServerResponse(ctx, http.StatusBadRequest, "Invalid payload")
			return
		}
	}
	resource.UserID = userID
	resource.Name = folderInput.Name
	resource.ResourceType = config.RESOURCE_TYPE_FOLDER
	resource.ParentID = folderInput.ParentID
	resource.FileExt = config.FOLDER_EXT
	resource.AccessType = config.ACCESS_TYPE_PRIVATE
	newResource, err := resource.CreateResource(db)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithData(ctx, http.StatusCreated, newResource.PublicResource())
}

func UpdateResource(ctx *gin.Context) {
	var resourceInput models.UpdateResourceInput
	db := config.NewDB()
	if err := ctx.ShouldBindJSON(&resourceInput); err != nil {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	resource := models.Resource{}
	resource.ID = resourceInput.ParentID
	if resourceInput.ParentID != "" {
		_, err := resource.FindResourceByID(db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ServerResponse(ctx, http.StatusNotFound, "parent_id not found")
			return
		}
	}
	fmt.Println(resourceInput)
	utils.SuccessWithMessage(ctx, http.StatusOK, "Updated successfully")
}
