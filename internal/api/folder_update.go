package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/gin-gonic/gin"
)



// UpdateFolderRoot updates folder root on database
//
// POST /api/folder/update/root
// formData: form.UpdateFolder
// @return 200 {string} string "ok"
func UpdateFolderRoot(router *gin.RouterGroup) {
	router.PUT("/folder/update/root", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)

		var form form.UpdateFolder

		if err := ctx.BindJSON(&form); err != nil {
			AbortBadRequest(ctx)
			return
		}

		folderMutex.Lock()
		defer folderMutex.Unlock()


		// Check if the current user is the owner of the folder
		isOwner, err := entity.IsFolderOwner(form.Uid, authPayload.UserID)
		if err != nil {
			AbortBadRequest(ctx)
			return
		}
		if !isOwner {
			// Forbidden: The user is not the owner of the folder
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this folder"})
			return
		}


		folder := entity.Folder{
			UID: form.Uid,
			Root:  form.Root,
		}

		if err := folder.UpdateRootOnly(); err != nil {
			AbortBadRequest(ctx)
			return
		}

		folder_user := entity.FolderUser{
			FolderID:   folder.ID,
			UserID:     authPayload.UserID,
			Permission: entity.OwnerPermission,
		}

		if err := folder_user.Update(); err != nil {
			AbortBadRequest(ctx)
			return
		}

		ctx.JSON(http.StatusOK, folder)
	})
}