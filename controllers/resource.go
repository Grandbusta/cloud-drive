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
		_, err := resource.FindResourceByID(db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ServerResponse(ctx, http.StatusNotFound, "parentId not found")
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
	utils.SuccessWithData(ctx, http.StatusCreated, newResource)
}
