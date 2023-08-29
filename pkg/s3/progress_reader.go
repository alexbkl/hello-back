package s3

import (
	"log"
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

	// I have no idea why the read length need to be div 2,
	// maybe the request read once when Sign and actually send call ReadAt again
	// It works for me
	log.Printf(
		"total read:%d    progress:%d%%\n",
		pr.read/2,
		int(float32(pr.read*100/2)/float32(pr.size)),
	)

	v := rds.UploadProgressValue{
		FileName: pr.file.Filename,
		Size:     pr.size,
		Read:     int64(float32(pr.read / 2)),
	}
	go pr.cb(pr.key, v)

	return n, err
}

func (pr *progressReader) Seek(offset int64, whence int) (int64, error) {
	return pr.src.Seek(offset, whence)
}
