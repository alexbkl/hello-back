package api

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/pkg/fs"
	"github.com/Hello-Storage/hello-back/pkg/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-gonic/gin"
)

// UploadFiles upload files to filebase using s3
//
// POST /api/file/upload
// Form: MultipartForm
// - files
func UploadFiles(router *gin.RouterGroup) {
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["files"]

		for _, file := range files {
			log.Infof("api: upload %s", file.Filename)

			mime, err := fs.GetFileContentType(file)

			if err != nil {
				log.Errorf("api: upload %s", err)
				AbortInternalServerError(c)
				return
			}

			f := entity.File{
				FileName: file.Filename,
				FileRoot: "/",
				FileMime: mime,
				FileSize: file.Size,
			}

			if err := f.Create(); err != nil {
				log.Errorf("api: upload %s", err)
				AbortInternalServerError(c)
				return
			}

			// if err := UploadFile(file); err != nil {
			// 	AbortInternalServerError(c)
			// 	return
			// }
		}
		c.JSON(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

	})
}

// internal upload one file
func UploadFile(file *multipart.FileHeader) error {
	config, _ := config.LoadEnv()

	s3Config := aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.FilebaseAccessKey, config.FilebaseSecretKey, ""),
		Endpoint:         aws.String("https://s3.filebase.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}

	err := s3.UploadObject(s3Config, file)

	return err
}
