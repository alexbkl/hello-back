package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/gin-gonic/gin"
)

// GetFile returns file details as JSON.
//
// GET /api/file/:hash
// Params:
// - hash (string) SHA-1 hash of the file
func GetFile(router *gin.RouterGroup) {
	router.GET("/file/:hash", func(c *gin.Context) {
		// To Do check access grant

		p, err := query.FileByHash(c.Param("hash"))

		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		c.JSON(http.StatusOK, p)
	})
}
