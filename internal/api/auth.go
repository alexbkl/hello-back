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
		var f form.LoginUserRequest
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
		var f form.RegisterUserRequest
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

// RequestNonce
// POST /api/nonce
func RequestNonce(router *gin.RouterGroup) {
	router.POST("/nonce", func(ctx *gin.Context) {
		var req struct {
			WalletAddress string `json:"wallet_address" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		u := entity.User{
			Wallet: entity.Wallet{
				Address: req.WalletAddress,
			},
		}

		nonce, err := u.RetrieveNonce(true)
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				ErrorResponse(err),
			)
			return
		}
		ctx.JSON(http.StatusOK, nonce)
	})
}
