package s3

import (
	"mime/multipart"
	"sync/atomic"

	"github.com/Hello-Storage/hello-back/internal/rds"
)

type progressReader struct {
	file *multipart.FileHeader
	src  multipart.File
	size int64
	read int64
	key  string
	cb   func(key string, val rds.UploadProgressValue)
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	return pr.src.Read(p)
}

func (pr *progressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := pr.src.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	// Got the length have read( or means has uploaded), and you can construct your message
	atomic.AddInt64(&pr.read, int64(n))

	v := rds.UploadProgressValue{
		Name: pr.file.Filename,
		Size: pr.size,
		Read: int64(float32(pr.read / 2)),
	}
	go pr.cb(pr.key, v)

	return n, err
}

func (pr *progressReader) Seek(offset int64, whence int) (int64, error) {
	return pr.src.Seek(offset, whence)
}
