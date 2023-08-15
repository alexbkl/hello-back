package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/gin-gonic/gin"
)

// RenewAccessToken
//
// POST /api/token/renew
func RenewAccessToken(router *gin.RouterGroup, tokenMaker token.Maker) {
	router.POST("/token/renew", func(ctx *gin.Context) {
		var req form.RenewAccessTokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		refreshPayload, err := tokenMaker.VerifyToken(req.RefreshToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		// TO-DO check session

		accessToken, accessPayload, err := tokenMaker.CreateToken(
			refreshPayload.UID,
			refreshPayload.Username,
			config.Env().AccessTokenDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		rsp := form.RenewAccessTokenResponse{
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessPayload.ExpiredAt,
		}
		ctx.JSON(http.StatusOK, rsp)
	})
}
