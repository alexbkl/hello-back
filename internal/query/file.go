package query

import (
	"fmt"

	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
)

// FileByUID finds a file entity for the given UID.
func FileByUID(fileUID string) (*entity.File, error) {
	f := entity.File{}

	if fileUID == "" {
		return &f, fmt.Errorf("file uid required")
	}

	err := db.Db().Where("file_uid = ?", fileUID).First(&f).Error

	return &f, err
}

// FileByHash finds a file with a given hash string.
func FileByHash(fileHash string) (*entity.File, error) {
	f := entity.File{}

	if fileHash == "" {
		return &f, fmt.Errorf("file hash required")
	}

	err := db.Db().Where("file_hash = ?", fileHash).First(&f).Error

	return &f, err
}
