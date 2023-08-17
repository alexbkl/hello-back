package s3

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadObject(s3Config aws.Config, file *multipart.FileHeader, bucket, key string) error {

	// create a new session using the config above and profile
	goSession, err := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: "filebase",
	})

	// check if the session was created correctly.
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// create a s3 client session
	s3Client := s3.New(goSession)

	// create put object input
	putObjectInput := &s3.PutObjectInput{
		Body:   src,
		Bucket: aws.String(bucket), // bucket name
		Key:    aws.String(fmt.Sprintf("%s", key)),
	}

	// upload file
	_, err = s3Client.PutObject(putObjectInput)

	if err != nil {
		return err
	}

	return nil
}
