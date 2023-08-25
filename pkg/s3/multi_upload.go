package s3

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadObjectV2(
	s3Config aws.Config,
	file *multipart.FileHeader,
	bucket, key string,
) error {

	// create a new session using the config above and profile
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: "wasabi",
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

	// Create an S3 client
	svc := s3.New(sess)

	// Step 1: Create a multipart upload
	createResp, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	// Get the upload ID for subsequent operations
	uploadID := createResp.UploadId

	// Set the part size and buffer for reading the file
	partSize := int64(5 * 1024 * 1024) // 5 MB
	buffer := make([]byte, partSize)

	// Store the uploaded part information
	var wg sync.WaitGroup
	var parts []*s3.CompletedPart
	var mu sync.Mutex

	// Read and upload each part
	for partNumber := int64(1); ; partNumber++ {
		numBytes, err := src.Read(buffer)
		if err != nil {
			break // End of file
		}

		wg.Add(1)

		// Perform each part upload concurrently
		go func(partNumber int64, content []byte) {
			defer wg.Done()

			uploadResp, err := svc.UploadPart(&s3.UploadPartInput{
				Bucket:        aws.String(bucket),
				Key:           aws.String(key),
				UploadId:      uploadID,
				PartNumber:    aws.Int64(partNumber),
				ContentLength: aws.Int64(int64(numBytes)),
				Body:          bytes.NewReader(content),
			})

			fmt.Println("partNumber:", partNumber)

			if err != nil {
				log.Printf("Failed to upload part %d: %v\n", partNumber, err)
				return
			}

			mu.Lock()
			parts = append(parts, &s3.CompletedPart{
				ETag:       uploadResp.ETag,
				PartNumber: aws.Int64(partNumber),
			})
			mu.Unlock()
		}(partNumber, buffer[:numBytes])

	}

	wg.Wait()

	// Step 2: Complete the multipart upload
	completeResp, err := svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		return err
	}

	fmt.Println("Upload completed. Location:", completeResp.Location)

	return nil
}
