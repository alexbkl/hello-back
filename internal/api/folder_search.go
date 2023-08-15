package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/gin-gonic/gin"
)

// FolderResponse represents the folder API response.
type FolderResponse struct {
	Root    string          `json:"root"`
	Folders []entity.Folder `json:"folders"`
	Files   []entity.File   `json:"files"`
}

// SearchFolders returns folders & files under request foldera s JSON.
//
// GET /api/folder/:uid
// Params:
// - uid FolderUID
func SearchFolderByRoot(router *gin.RouterGroup) {

	handler := func(ctx *gin.Context, root string) {
		// TO-DO check access grant

		resp := FolderResponse{Root: root}

		// TO-DO folders
		if folders, err := query.FoldersByRoot(root); err != nil {
			log.Errorf("folder: %s", err)

			AbortInternalServerError(ctx)
			return
		} else {
			resp.Folders = folders
		}

		// files
		if files, err := query.FilesByRoot(root); err != nil {
			log.Errorf("file: %s", err)

			AbortInternalServerError(ctx)
			return
		} else {
			resp.Files = files
		}

		ctx.JSON(http.StatusOK, resp)
	}

	router.GET("/folder", func(ctx *gin.Context) {
		handler(ctx, "/")
	})

	router.GET("/folder/:uid", func(ctx *gin.Context) {
		uid := ctx.Param("uid")

		handler(ctx, uid)
	})
}
