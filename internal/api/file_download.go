package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/Hello-Storage/hello-back/pkg/s3"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// DownloadFile downloads file from filebase using s3
//
// GET /api/file/download/:uid
// @param uid path string true "file uid"
// @return 200 {string} string "ok"
func DownloadFile(router *gin.RouterGroup) {
	router.GET("/download/:uid", func(ctx *gin.Context) {
		// TO-DO check user auth & add user uid
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)

		u := query.FindUser(entity.User{ID: authPayload.UserID})
		key := ctx.Param("uid")
		log.Infof("u : %v", authPayload.UserID)
		// Multipart form
		keyPath := u.UID + "/" + key
		out, error := DownloadFileFromS3(keyPath)

		if error != nil {
			log.Errorf("download file: %s", error)
			//log key:
			log.Infof("key : %v", key)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": error.Error(),
			})
			return
		}

		// fix u declared and not used error
		log.Infof("u : %v", u)

		// Set the correct content type and file name
		ctx.Header("Content-Type", *out.ContentType)
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", key))

		// Copy the file data to the response
		_, error = io.Copy(ctx.Writer, out.Body)
		if error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"result":  out,
		})
	})
}

// internal upload one file
func DownloadFileFromS3(key string) (*awsS3.GetObjectOutput, error) {

	s3Config := aws.Config{
		Credentials: credentials.NewStaticCredentials(
			config.Env().WasabiAccessKey,
			config.Env().WasabiSecretKey,
			"",
		),
		Endpoint:         aws.String(config.Env().WasabiEndpoint),
		Region:           aws.String(config.Env().WasabiRegion),
		S3ForcePathStyle: aws.Bool(true),
	}

	out, err := s3.DownloadObject(s3Config, key)

	return out, err
}
