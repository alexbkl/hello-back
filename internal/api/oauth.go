package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/form"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/Hello-Storage/hello-back/pkg/oauth"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/gin-gonic/gin"
)

type OAuthPayload struct {
	Token string `json:"token" binding:"required"`
}

// OAuthGoogle
//
// POST /api/oauth/google
func OAuthGoogle(router *gin.RouterGroup, tokenMaker token.Maker) {
	router.POST("/oauth/google", func(ctx *gin.Context) {
		var f OAuthPayload

		if err := ctx.ShouldBindJSON(&f); err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		google_user, err := oauth.GetGoogleUser(f.Token)

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		u := query.FindUser(entity.User{Email: google_user.Email})
		if u == nil {
			// create new user
			new := entity.User{
				Name:  google_user.Name,
				Email: google_user.Email,
			}

			if err := new.Save(); err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}

			u = &new
		}

		// authorization token
		accessToken, accessPayload, err := tokenMaker.CreateToken(
			u.Name,
			config.Env().AccessTokenDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		refreshToken, refreshPayload, err := tokenMaker.CreateToken(
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

// OAuthGithub
//
// GET /api/oauth/github
func OAuthGithub(router *gin.RouterGroup, tokenMaker token.Maker) {
	router.POST("/oauth/github", func(ctx *gin.Context) {
		// var f OAuthPayload

		// code := ctx.Query("code")
		// state := ctx.Query("state")

		// if err := ctx.ShouldBindJSON(&f); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		// 	return
		// }
	})
}
