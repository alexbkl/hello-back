package api

import (
	"net/http"
	"sync"

	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/gin-gonic/gin"
)

var folderMutex = sync.Mutex{}

func CreateFolder(router *gin.RouterGroup) {
	router.POST("/create", func(ctx *gin.Context) {
		var f form.CreateFolder

		if err := ctx.BindJSON(&f); err != nil {
			AbortBadRequest(ctx)
			return
		}

		folderMutex.Lock()
		defer folderMutex.Unlock()

		m := entity.Folder{
			Title: f.Title,
		}

		if err := m.Create(); err != nil {
			AbortBadRequest(ctx)
			return
		}

		ctx.JSON(http.StatusOK, m)
	})
}
