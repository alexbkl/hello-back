package api

import "github.com/gin-gonic/gin"

// UpdateUser updates the profile information of the currently authenticated user.
//
// PUT /api/v1/users/:uid
func UpdateUser(router *gin.RouterGroup) {
	router.PUT("/user/:id", func(ctx *gin.Context) {

	})
}
