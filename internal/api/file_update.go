package api

import (
	"net/http"
	"sync"

	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/gin-gonic/gin"
)

var fileMutex = sync.Mutex{}

// UpadteFileRoot updates file root on database
//
// POST /api/file/update/root
// formData: form.UpdateFile
// @return 200 {string} string "ok"
func UpdateFileRoot(router *gin.RouterGroup) {
	router.PUT("/file/update/root", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)

		var form form.UpdateFileRoot

		if err := ctx.BindJSON(&form); err != nil {
			AbortBadRequest(ctx)
			return
		}

		fileMutex.Lock()
		defer fileMutex.Unlock()

		// Check if the current user is the owner of the folder
		isOwner, err := entity.IsFileOwner(form.Uid, authPayload.UserID)
		if err != nil {
			AbortBadRequest(ctx)
			return
		}
		if !isOwner {
			// Forbidden: The user is not the owner of the folder
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this file"})
			return
		}

		file := entity.File{}

		// Query the database to populate the 'File' entity with existing data.
		if err := db.Db().Where("UID = ?", form.Uid).First(&file).Error; err != nil {
			AbortBadRequest(ctx)
			return
		}
		file.UID = form.Uid
		file.Root = form.Root

		if err := file.UpdateRootOnly(); err != nil {
			AbortBadRequest(ctx)
			return
		}

		file_user := entity.FileUser{
			FileID:     file.ID,
			UserID:     authPayload.UserID,
			Permission: entity.OwnerPermission,
		}

		if err := file_user.Update(); err != nil {
			AbortBadRequest(ctx)
			return
		}

		ctx.JSON(http.StatusOK, file)
	})
}
