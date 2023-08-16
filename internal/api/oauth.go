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

// OAuthGoogle
//
// GET /api/oauth/google
func OAuthGoogle(router *gin.RouterGroup, tokenMaker token.Maker) {
	router.GET("/oauth/google", func(ctx *gin.Context) {

		code := ctx.Query("code")

		if code == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
			return
		}

		google_user, err := oauth.GetGoogleUser(code)

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		u := query.FindUserByEmail(google_user.Email)
		if u == nil {
			// create new user
			new := entity.User{
				Name: google_user.Name,
				Email: entity.Email{
					Email: google_user.Email,
				},
			}

			if err := new.Create(); err != nil {
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}

			u = &new
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

// OAuthGithub
//
// GET /api/oauth/github
func OAuthGithub(router *gin.RouterGroup, tokenMaker token.Maker) {
	router.GET("/oauth/github", func(ctx *gin.Context) {
		code := ctx.Query("code")

		if code == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
			return
		}

		token, err := oauth.GetGithubOAuthToken(code)

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		github_user, err := oauth.GetGithubUser(token)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		u := query.FindUserByGithub(github_user.ID)
		if u == nil {
			new := entity.User{
				Name: github_user.Name,
				Github: entity.Github{
					GithubID: github_user.ID,
					Name:     github_user.Name,
					Avatar:   github_user.Avatar,
				},
			}

			if err := new.Create(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
				return
			}

			u = &new
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
