package api

import (
	"fmt"

	"github.com/Hello-Storage/hello-back/internal/config"
	"github.com/Hello-Storage/hello-back/internal/constant"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/query"
	"github.com/Hello-Storage/hello-back/pkg/s3"
	"github.com/Hello-Storage/hello-back/pkg/token"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-gonic/gin"
)

// DownloadFile download file from filebase s3
//
// DELETE /api/file/delete/:uid
//
// @param uid path string true "file uid"
// @return 200 {string} string "ok"

func DeleteFile(router *gin.RouterGroup) {
	router.DELETE("/delete/:uid", func(ctx *gin.Context) {
		// TO-DO check user auth & add user uid
		authPayload := ctx.MustGet(constant.AuthorizationPayloadKey).(*token.Payload)

		file_uid := ctx.Param("uid")

		f, err := query.FindFileByUID(file_uid)

		if err != nil {
			AbortEntityNotFound(ctx)
			log.Errorf("file not found: %v", err)
			return
		}

		f_u := entity.FileUser{
			FileID: f.ID,
			UserID: authPayload.UserID,
		}

		//delete file from s3
		if err := DeleteFileFromS3(fmt.Sprintf("%s/%s", authPayload.UserUID, file_uid)); err != nil {
			AbortInternalServerError(ctx)
			log.Errorf("delete file from s3 error: %v", err)
			return
		}

		//Delete file
		if err := query.DeleteFileByUID(f.UID); err != nil {
			AbortInternalServerError(ctx)
			log.Errorf("delete file error: %v", err)
			return
		}

		//Delete file user
		if err := query.DeleteFileUser(f_u); err != nil {
			AbortInternalServerError(ctx)
			log.Errorf("delete file user error: %v", err)
			return
		}

		// remove user storage quantity
		user_detail := query.FindUserDetailByUserID(authPayload.UserID)

		if err := user_detail.Update("storage_used", user_detail.StorageUsed-uint(f.Size)); err != nil {
			log.Errorf("removing storage_used: %s", err)
			AbortInternalServerError(ctx)
			return
		}

		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})
}

// internal delete one file
func DeleteFileFromS3(key string) error {
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

	//delete file from s3
	if err := s3.DeleteObject(s3Config, config.Env().FilebaseBucket, key); err != nil {
		log.Errorf("DeleteFileFromS3: delete file from s3 error: %v", err)
		return err
	}

	return nil
}
