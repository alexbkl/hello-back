package s3client

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

var S3Client *s3.S3

func UploadFile() {

	FILEBASE_ACCESS_KEY := os.Getenv("FILEBASE_ACCESS_KEY")
	FILEBASE_SECRET_ACCESS_KEY := os.Getenv("FILEBASE_SECRET_ACCESS_KEY")

	bucketName := "hello-storage"
	filename := "example.txt"
	region := "us-east-1"
	endpoint := "https://s3.filebase.com"

    
	///configuration
	s3Config := aws.Config{
		Credentials:      credentials.NewStaticCredentials(FILEBASE_ACCESS_KEY, FILEBASE_SECRET_ACCESS_KEY, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}

	goSession, err := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: "filebase",
	})

	// check if the session was created correctly.

	if err != nil {
		fmt.Println(err)
	}

	// create a s3 client session
	s3Client := s3.New(goSession)

    //assign the s3Client to the global variable

	// set parameter for bucket name
	bucket := aws.String(bucketName)

	// create a bucket
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})

	// print if there is an error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
    S3Client = s3Client

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	
	// create put object input
	putObjectInput := &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}

	//upload file
	_, err = s3Client.PutObject(putObjectInput)
	// print if there is an error
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("File uploaded from " + filename + " to " + bucketName)
	}


}
