package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/gin-gonic/gin"
)

// GetFile returns file details as JSON.
//
// GET /api/file/info/:uid
// Params:
// - uid
func GetFile(router *gin.RouterGroup) {
	router.GET("/info/:uid", func(c *gin.Context) {
		// To Do check access grant
		uid := c.Param("uid")

		p, err := query.FindFileByUID(uid)

		if err != nil {
			AbortEntityNotFound(c)
			return
		}

		c.JSON(http.StatusOK, p)
	})
}
