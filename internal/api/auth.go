package api

import (
	"net/http"
	"sync"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/Hello-Storage/hello-back/pkg/web3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var authMutex = sync.Mutex{}

// LoadUser
//
// GET /api/load
func LoadUser(router *gin.RouterGroup) {
	router.GET("/load", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)

		u := query.FindUser(entity.User{Model: gorm.Model{ID: authPayload.UserID}})
		if u == nil {
			ctx.JSON(http.StatusNotFound, "user not found")
			return
		}

		var resp = struct {
			UID           string `json:"uid"`
			Name          string `json:"name"`
			Role          string `json:"role"`
			WalletAddress string `json:"walletAddress"`
		}{
			UID:           u.UID,
			Name:          u.Name,
			Role:          string(u.Role),
			WalletAddress: u.Wallet.Address,
		}

		ctx.JSON(http.StatusOK, resp)
	})
}

// LoginUser
//
// POST /api/login
func LoginUser(router *gin.RouterGroup, tokenMaker token.Maker) {
	router.POST("/login", func(ctx *gin.Context) {
		var f form.LoginUserRequest
		if err := ctx.BindJSON(&f); err != nil {
			AbortBadRequest(ctx)
			return
		}

		authMutex.Lock()
		defer authMutex.Unlock()

		u := query.FindUserByWalletAddress(f.WalletAddress)
		if u == nil {
			Abort(ctx, http.StatusNotFound, "user not exists!")
			return
		}

		// retrieve nonce
		nonce, err := u.RetrieveNonce(false)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		log.Infof("nonce: %s", nonce)

		// validate signature
		result := web3.ValidateMessageSignature(f.WalletAddress, f.Signature, constant.BuildLoginMessage(nonce))
		if !result {
			ctx.JSON(http.StatusBadRequest, "invalide signature")
			return
		}

		// authorization token
		accessToken, accessPayload, err := tokenMaker.CreateToken(
			u.ID,
			u.UID,
			u.Name,
			config.Env().AccessTokenDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		refreshToken, refreshPayload, err := tokenMaker.CreateToken(
			u.ID,
			u.UID,
			u.Name,
			config.Env().RefreshTokenDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		// TO-DO create session part

		rsp := form.LoginUserResponse{
			// SessionID:             session.ID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		}
		ctx.JSON(http.StatusOK, rsp)
	})

}

// RegisterUser
//
// POST /api/register
func RegisterUser(router *gin.RouterGroup, tokenMaker token.Maker) {
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

		log.Info("renew", u)
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
