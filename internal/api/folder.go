package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/gin-gonic/gin"
)

// GetFolder returns folders & files under request foldera s JSON.
//
// GET /api/folder/:uid
// Params:
// - uid FolderUID
func GetFolder(router *gin.RouterGroup) {
	router.GET("/folder/:uid", func(c *gin.Context) {
		// To Do check access grant

		p, err := query.FileByHash(c.Param("uid"))

		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		c.JSON(http.StatusOK, p)
	})
}
