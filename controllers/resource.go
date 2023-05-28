package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Grandbusta/cloud-drive/config"
	"github.com/Grandbusta/cloud-drive/models"
	"github.com/Grandbusta/cloud-drive/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateFolder(ctx *gin.Context) {
	db := config.NewDB()
	var folderInput models.CreateFolderInput
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
	parent := models.Resource{}
	parent.ID = folderInput.ParentID
	if folderInput.ParentID != config.ROOT {
		parentResource, err := parent.FindResourceByID(db)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ServerResponse(ctx, http.StatusNotFound, "parent_id not found")
			return
		}
		if parentResource.ResourceType != config.RESOURCE_TYPE_FOLDER {
			utils.ServerResponse(ctx, http.StatusBadRequest, "Invalid payload")
			return
		}
		// resource.Path = parentResource.Path
	}
	resource.ParentID = folderInput.ParentID
	resource.UserID = userID
	resource.Name = folderInput.Name
	resource.ResourceType = config.RESOURCE_TYPE_FOLDER
	resource.FileExt = config.FOLDER_EXT
	resource.AccessType = config.ACCESS_TYPE_PRIVATE
	newResource, err := resource.CreateResource(db)
	treePath := models.TreePath{}
	if folderInput.ParentID == config.ROOT {
		treePath.Ancestor = newResource.ID
		treePath.Descendant = newResource.ID
		err = treePath.InsertRoot(db)
	} else {
		treePath.Ancestor = newResource.ParentID
		treePath.Descendant = newResource.ID
		err = treePath.InsertDescendant(db)
	}
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithData(ctx, http.StatusCreated, newResource.PublicResource())
}

func UpdateResource(ctx *gin.Context) {
	db := config.NewDB()
	var resourceInput models.UpdateResourceInput
	resource_id := ctx.Param("resource_id")
	if err := ctx.ShouldBindJSON(&resourceInput); err != nil {
		utils.ServerResponse(ctx, http.StatusUnprocessableEntity, "Invalid payload")
		return
	}
	resource := models.Resource{}
	parent := models.Resource{}
	resource.ID = resource_id
	parent.ID = resourceInput.ParentID

	_, err := resource.FindResourceByID(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ServerResponse(ctx, http.StatusNotFound, "resource not found")
		return
	}

	if resourceInput.ParentID != "" {
		if resourceInput.ParentID != config.ROOT {
			parentResource, err := parent.FindResourceByID(db)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.ServerResponse(ctx, http.StatusNotFound, "parent_id not found")
				return
			}
			if parentResource.ResourceType != config.RESOURCE_TYPE_FOLDER {
				utils.ServerResponse(ctx, http.StatusBadRequest, "Invalid payload")
				return
			}
			// resource.Path = parentResource.Path
			resource.ParentID = parentResource.ID
		}
	}

	if resourceInput.Name != "" {
		resource.Name = resourceInput.Name
	}

	if resourceInput.AccessType != "" {
		if resourceInput.AccessType == config.ACCESS_TYPE_PRIVATE || resourceInput.AccessType == config.ACCESS_TYPE_PUBLIC {
			resource.AccessType = resourceInput.AccessType
		} else {
			utils.ServerResponse(ctx, http.StatusBadRequest, "Invalid payload")
			return
		}
	}

	_, err = resource.UpdateResource(db)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}

	utils.SuccessWithMessage(ctx, http.StatusOK, "Updated successfully")
}

func DeleteResource(ctx *gin.Context) {
	db := config.NewDB()
	resource_id := ctx.Param("resource_id")
	resource := models.Resource{}
	resource.ID = resource_id
	_, err := resource.FindResourceByID(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ServerResponse(ctx, http.StatusNotFound, "resource not found")
		return
	}
	err = resource.DeleteResource(db)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithMessage(ctx, http.StatusOK, "Deleted successfully")
}

func UploadFile(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	log.Println(file.Filename, file.Header)
}
