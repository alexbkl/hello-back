package s3client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"meta-go-api/config"
	"meta-go-api/entities"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ipfs/go-cid"
)

var S3Client *s3.S3
var BucketName string = "hello-storage"
var Region string = "us-east-1"
var Endpoint string = "https://s3.filebase.com"
var Bucket *string = aws.String(BucketName)

var FILEBASE_ACCESS_KEY string
var FILEBASE_SECRET_ACCESS_KEY string

func Init() {

	//set environment variables
	FILEBASE_ACCESS_KEY = os.Getenv("FILEBASE_ACCESS_KEY")
	FILEBASE_SECRET_ACCESS_KEY = os.Getenv("FILEBASE_SECRET_ACCESS_KEY")

	///configuration
	s3Config := aws.Config{
		Credentials:      credentials.NewStaticCredentials(FILEBASE_ACCESS_KEY, FILEBASE_SECRET_ACCESS_KEY, ""),
		Endpoint:         aws.String(Endpoint),
		Region:           aws.String(Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	goSession, err := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: "filebase",
	})

	// check if the session was created correctly.

	if err != nil {
		fmt.Println("Error creating session ", err)
	}

	// create a s3 client session
	s3Client := s3.New(goSession)

	//assign the s3Client to the global variable

	// set parameter for bucket name

	bucket := aws.String("hello-storage")

	// create a bucket
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})

	// print if there is an error
	if err != nil {
		fmt.Println("Error creating bucket ", err)
	}

	S3Client = s3Client
}


func UploadFile(cidOfEncryptedBufferStr string, src multipart.File) ([]byte, error) {

	bucketName := "hello-storage"

	//transform src to []byte
	srcBytes, err := ioutil.ReadAll(src)
	


	//if cid exists in database and name of the file is the same, return the file without uploading to s3
	var fileExists entities.File

	//check if cid exists in database
	result := config.Database.Where("c_id_of_encrypted_buffer = ?", cidOfEncryptedBufferStr).Find(&fileExists)
	//print the result
	if result.RowsAffected != 0 {
		//fmt.Print("File already exists")		
		return srcBytes, nil
	}

	/*
		// Create a cid from a marshaled string
		decodedC, err := cid.Decode(c.String())
		if err != nil {
			fmt.Println("Error decoding CID: ", err)
			return nil, err
		}
		fmt.Println("Got CID: ", decodedC)
	*/
	/*
	metadata := map[string]*string{
		"Content-Type":      aws.String(fileHeader.Header.Get("Content-Type")),
		"Original-Filename": aws.String(fileHeader.Filename),
		"Content-Length":    aws.String(fmt.Sprintf("%d", fileHeader.Size)),
	}*/
	// create put object input
	putObjectInput := &s3.PutObjectInput{
		Body:     bytes.NewReader(srcBytes),
		Bucket:   aws.String(bucketName),
		Key:      aws.String(cidOfEncryptedBufferStr),
		//Metadata: metadata,
	}
	//upload file
	_, err = S3Client.PutObject(putObjectInput)
	// print if there is an error
	if err != nil {
		fmt.Println("Error uploading file ", err)
		return nil, err
	}
	//err = PinCID(c)

	return srcBytes, nil

}

func DownloadFile(cid string) (*s3.GetObjectOutput, error) {
	// create put object input
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(cid),
	}

	// get file
	//result is *s3.GetObjectOutput type
	result, err := S3Client.GetObject(getObjectInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Metadata:")
	for key, value := range result.Metadata {
		// Need to dereference the value pointer to get the actual string.
		fmt.Printf("  %s: %s\n", key, *value)
	}

	return result, nil
}

func PinCID(c cid.Cid) error {
	accessToken := os.Getenv("FILEBASE_PINNING_ACCESS_TOKEN")

	url := "https://api.filebase.io/v1/ipfs/pins"
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Content-Type":  "application/json",
	}
	body := fmt.Sprintf(`{"cid": "%s"}`, c.String())

	response, err := SendRequest("POST", url, headers, body)
	if err != nil {
		fmt.Println("Error pinning CID: ", err)
		return err
	}

	fmt.Println("Pinned CID: ", c, "Response: ", response)

	return nil
}

func DeleteFile(cid string) error {
	// create put object input
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: Bucket,
		Key:    aws.String(cid),
	}

	// get file
	_, err := S3Client.DeleteObject(deleteObjectInput)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
