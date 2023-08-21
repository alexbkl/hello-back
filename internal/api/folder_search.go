package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/gin-gonic/gin"
)

// FolderResponse represents the folder API response.
type FolderResponse struct {
	Root    string         `json:"root"`
	Path    entity.Folders `json:"path"`
	Folders entity.Folders `json:"folders"`
	Files   entity.Files   `json:"files"`
}

// SearchFolders returns folders & files under request foldera s JSON.
//
// GET /api/folder/:uid
// Params:
// - uid FolderUID
func SearchFolderByRoot(router *gin.RouterGroup) {

	handler := func(ctx *gin.Context, root string) {
		// TO-DO check access grant
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)
		resp := FolderResponse{Root: root}

		// TO-DO folders
		if root == "/" {
			if folders, err := query.FindRootFoldersByUser(authPayload.UserID); err != nil {
				log.Errorf("folder: %s", err)

				AbortInternalServerError(ctx)
				return
			} else {
				resp.Folders = folders
			}

			// files
			if files, err := query.FindRootFilesByUser(authPayload.UserID); err != nil {
				log.Errorf("file: %s", err)

				AbortInternalServerError(ctx)
				return
			} else {
				resp.Files = files
			}

		} else {
			if folders, err := query.FoldersByRoot(root); err != nil {
				log.Errorf("folder: %s", err)

				AbortInternalServerError(ctx)
				return
			} else {
				resp.Folders = folders
			}

			// files
			if files, err := query.FindFilesByRoot(root); err != nil {
				log.Errorf("file: %s", err)

				AbortInternalServerError(ctx)
				return
			} else {
				resp.Files = files
			}
		}

		// add path
		resp.Path = query.FindFolderPathByRoot(root)

		ctx.JSON(http.StatusOK, resp)
	}

	router.GET("/folder", func(ctx *gin.Context) {
		handler(ctx, "/")
	})

	router.GET("/folder/:uid", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)
		uid := ctx.Param("uid")

		// check if user own this root
		if perm := query.CheckFolderPermByUser(uid, authPayload.UserID); !perm {
			ctx.JSON(http.StatusForbidden, "don't have access")
			return
		}

		handler(ctx, uid)
	})
}
