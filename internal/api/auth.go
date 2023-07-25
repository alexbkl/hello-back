package api

import (
	"net/http"
	"sync"

	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/gin-gonic/gin"
)

var authMutex = sync.Mutex{}

// LoginUser
//
// POST /api/login
func LoginUser(router *gin.RouterGroup) {
	router.POST("/login", func(ctx *gin.Context) {
		var f form.User
		if err := ctx.BindJSON(&f); err != nil {
			AbortBadRequest(ctx)
			return
		}

		authMutex.Lock()
		defer authMutex.Unlock()

		u := entity.User{
			Name: f.Name,
		}

		// TO-DO check exists user info, if
		if user := query.FindUser(u); user != nil {
			Abort(ctx, http.StatusBadRequest, "user already exists!")
		}

		if err := u.Create(); err != nil {
			AbortInternalServerError(ctx)
			return
		}

		ctx.JSON(
			http.StatusOK,
			"user created!",
		)
	})

}

// RegisterUser
//
// POST /api/register
func RegisterUser(router *gin.RouterGroup) {
	router.POST("/register", func(ctx *gin.Context) {
		var f form.User
		if err := ctx.BindJSON(&f); err != nil {
			AbortBadRequest(ctx)
			return
		}

		authMutex.Lock()
		defer authMutex.Unlock()

		u := entity.User{
			Name: f.Name,
		}

		// TO-DO check exists user info, if
		if user := query.FindUser(u); user != nil {
			Abort(ctx, http.StatusBadRequest, "user already exists!")
		}

		if err := u.Create(); err != nil {
			AbortInternalServerError(ctx)
			return
		}

		ctx.JSON(
			http.StatusOK,
			"user created!",
		)
	})
}
