package api

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/pkg/fs"
	"github.com/Hello-Storage/hello-back/pkg/s3"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-gonic/gin"
)

// UploadFiles upload files to filebase using s3
//
// POST /api/file/upload
// Form: MultipartForm
// - files
// - root
func UploadFiles(router *gin.RouterGroup) {
	router.POST("/upload", func(ctx *gin.Context) {
		// TO-DO check user auth & add user uid
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)

		// Multipart form
		form, err := ctx.MultipartForm()

		if err != nil {
			AbortBadRequest(ctx)
			return
		}

		files := form.File["files"]
		root := form.Value["root"]

		var r string
		if len(root) > 0 {
			r = root[0]
			log.Infof("file: %s", r)
		} else {
			r = "/"
		}

		for _, file := range files {
			log.Infof("api: upload %s", file.Filename)

			mime, err := fs.GetFileContentType(file)

			if err != nil {
				log.Errorf("api: upload %s", err)
				AbortInternalServerError(ctx)
				return
			}

			f := entity.File{
				Name: file.Filename,
				Root: r,
				Mime: mime,
				Size: file.Size,
			}

			// upload file
			if err := UploadFile(file, authPayload.UID, f.UID); err != nil {
				log.Errorf("api: upload %s", err)
				AbortInternalServerError(ctx)
				return
			}

			// save file info to db
			if err := f.Create(); err != nil {
				log.Errorf("api: upload %s", err)
				AbortInternalServerError(ctx)
				return
			}
		}
		ctx.JSON(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
}

// internal upload one file
func UploadFile(file *multipart.FileHeader, user_uid, key string) error {

	s3Config := aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Env().FilebaseAccessKey, config.Env().FilebaseSecretKey, ""),
		Endpoint:         aws.String("https://s3.filebase.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}

	err := s3.UploadObject(s3Config, file, config.Env().FilebaseBucket, user_uid, key)

	return err
}
