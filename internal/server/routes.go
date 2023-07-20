package server

import (
	"github.com/Hello-Storage/hello-back/internal/api"
	"github.com/gin-gonic/gin"
)

var APIv1 *gin.RouterGroup

func registerRoutes(router *gin.Engine) {
	// Enables automatic redirection if the current route cannot be matched but a
	// handler for the path with (without) the trailing slash exists.
	// router.RedirectTrailingSlash = true

	// routes
	api.Ping(APIv1)

	api.UpdateUser(APIv1)
	api.GetFile(APIv1)

	api.UploadFiles(APIv1)
}
