package api

import (
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
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
		} else {
			r = "/"
		}

		for _, file := range files {
			_, params, err := mime.ParseMediaType(file.Header.Get("Content-Disposition"))
			if err != nil {
				log.Errorf("parse media type: %s", err)
				AbortInternalServerError(ctx)
				return
			}

			mime := file.Header.Get("Content-Type")

			file_path := params["filename"]
			actual_root, err := GetAndProcessFileRoot(file_path, r, authPayload.UserID)
			log.Infof("actual_root: %s", actual_root)

			f := entity.File{
				Name: file.Filename,
				Root: actual_root,
				Mime: mime,
				Size: file.Size,
			}
			if err := f.Create(); err != nil {
				AbortInternalServerError(ctx)
				return
			}

			// create file_user relation
			f_u := entity.FileUser{
				FileID: f.ID,
				UserID: authPayload.UserID,
			}
			if err := f_u.Create(); err != nil {
				AbortInternalServerError(ctx)
				return
			}

			// upload file
			if err := UploadFileToS3(file, f.UID); err != nil {
				log.Errorf("api: upload %s", err)
				AbortInternalServerError(ctx)
				return
			}

			// save file info to db

		}
		ctx.JSON(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
}

// internal upload one file
func UploadFileToS3(file *multipart.FileHeader, key string) error {

	s3Config := aws.Config{
		Credentials: credentials.NewStaticCredentials(
			config.Env().FilebaseAccessKey,
			config.Env().FilebaseSecretKey,
			"",
		),
		Endpoint:         aws.String("https://s3.filebase.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}

	err := s3.UploadObject(s3Config, file, config.Env().FilebaseBucket, key)

	return err
}

// internal
// here root => uid format
func GetAndProcessFileRoot(file_path, root string, user_id uint) (string, error) {
	res := strings.Split(file_path, "/")
	if len(res) == 1 {
		return root, nil
	}

	sub_file_path := strings.Join(res[1:], "/")
	sub_title := res[0]

	f := (&entity.Folder{
		Title: sub_title,
		Root:  root,
	}).FirstOrCreateFolderByTitleAndRoot()

	if f == nil {
		return "", errors.New("can't create folder")
	}

	// create folder_user relation
	f_u := &entity.FolderUser{
		FolderID:   f.ID,
		UserID:     user_id,
		Permission: entity.OwnerPermission,
	}

	if err := f_u.Create(); err != nil {
		return "", errors.New("can't create folder_user relation")
	}

	return GetAndProcessFileRoot(sub_file_path, f.UID, user_id)
}
