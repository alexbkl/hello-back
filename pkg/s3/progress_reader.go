package s3

import (
	"log"
	"mime/multipart"
	"sync/atomic"
)

type progressReader struct {
	fp   multipart.File
	size int64
	read int64
	key  string
	cb   func(key string, val int)
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	return pr.fp.Read(p)
}

func (pr *progressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := pr.fp.ReadAt(p, off)
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

	go pr.cb(pr.key, int(float32(pr.read*100/2)/float32(pr.size)))

	return n, err
}

func (pr *progressReader) Seek(offset int64, whence int) (int64, error) {
	return pr.fp.Seek(offset, whence)
}
