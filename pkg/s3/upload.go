package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadObject(s3Config aws.Config) error {
	// create a new session using the config above and profile
	goSession, err := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: "filebase",
	})

	// check if the session was created correctly.
	if err != nil {
		return err
	}

	// create a s3 client session
	s3Client := s3.New(goSession)

	//set the file path to upload
	file, err := os.Open("/path/to/object/to/upload")
	if err != nil {
		return err
	}

	defer file.Close()
	// create put object input
	putObjectInput := &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String("bucket-name"),
		Key:    aws.String("object-name"),
	}

	// upload file
	_, err = s3Client.PutObject(putObjectInput)

	if err != nil {
		return err
	}

	return nil
}
