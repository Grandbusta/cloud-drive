package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Grandbusta/cloud-drive/config"
	"github.com/Grandbusta/cloud-drive/models"
	"github.com/Grandbusta/cloud-drive/utils"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	resource.ID = resource_id

	_, err := resource.FindResourceByID(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ServerResponse(ctx, http.StatusNotFound, "resource not found")
		return
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
	treePath := models.TreePath{}
	treePath.Ancestor = resource_id
	resourceIds, err := treePath.SelectSelfWithDescendants(db)
	err = treePath.DeleteSelfWithDescendants(db)
	err = models.DeleteResourceByIds(db, resourceIds)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithMessage(ctx, http.StatusOK, "Deleted successfully")
}

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.ServerResponse(ctx, http.StatusBadRequest, "Failed to upload")
		return
	}
	log.Println(file.Filename, file.Header, file)
	cld := config.NewCld()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "cloud-drive", ResourceType: "auto"})
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithMessageAndData(ctx, http.StatusOK, "Upload successful", uploadResult)
}

func GetResource(ctx *gin.Context) {
	db := config.NewDB()
	resource_id := ctx.Param("resource_id")
	resource := models.Resource{}
	resource.ID = resource_id
	existingResource, err := resource.FindResourceByID(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ServerResponse(ctx, http.StatusNotFound, "resource not found")
		return
	}
	if existingResource.ResourceType == config.RESOURCE_TYPE_FOLDER {
		treePath := models.TreePath{}
		treePath.Ancestor = resource_id
		res, err := treePath.SelectChildren(db)
		if err != nil {
			utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
			return
		}
		if len(res) == 0 {
			res = make([]models.PublicResource, 0)
		}
		utils.SuccessWithData(ctx, http.StatusOK,
			map[string]interface{}{
				"resource": existingResource.PublicResource(),
				"children": res,
			},
		)
		return
	}
	utils.SuccessWithData(ctx, http.StatusOK, existingResource.PublicResource())

}
