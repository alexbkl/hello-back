package api

import (
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/gin-gonic/gin"
)

// UpdateUser updates the profile information of the currently authenticated user.
//
// PUT /api/v1/users/:uid
func UpdateUser(router *gin.RouterGroup) {
	router.GET("/user/:id", func(ctx *gin.Context) {
		user := entity.User{
			Name: "abc",
		}

		err := user.Create()

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		ctx.Status(200)
	})
}
