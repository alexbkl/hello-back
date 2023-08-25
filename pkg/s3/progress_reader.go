package s3

import (
	"fmt"
	"io"
	"log"
)

type progressReader struct {
	Reader     io.Reader
	TotalBytes int64
	BytesRead  int64
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.BytesRead += int64(n)
	log.Printf("Upload progress: %d/%d bytes", pr.BytesRead, pr.TotalBytes)
	return n, err
}

type ProgressListener struct {
	TotalBytes       int64
	BytesTransferred int64
	UploadID         string
}

func (pl *ProgressListener) OnPartUploadCompleted(partNum int, numBytes int64) {
	pl.BytesTransferred += numBytes
	fmt.Printf("Uploaded part %d: %d / %d bytes\n", partNum, pl.BytesTransferred, pl.TotalBytes)
}
