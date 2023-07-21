/*
Package api provides REST API authentication and request handlers.
*/

package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/event"
	"github.com/gin-gonic/gin"
)

var log = event.Log

func Ping(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
