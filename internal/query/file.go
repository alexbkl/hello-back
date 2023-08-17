package query

import (
	"fmt"

	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
)

// FileByUID returns file for the given UID.
func FileByUID(uid string) (*entity.File, error) {
	f := entity.File{}

	if uid == "" {
		return &f, fmt.Errorf("file uid required")
	}

	err := db.Db().Where("uid = ?", uid).First(&f).Error

	return &f, err
}

// FilesByRoot return files in a given folder root.
func FilesByRoot(root string) (files entity.Files, err error) {
	if err := db.Db().Where("root = ?", root).Find(&files).Error; err != nil {
		return files, err
	}

	return files, err
}

// DeleteFileByUID deletes a file by its UID.
func DeleteFileByUID(fileUid string) error {
	if fileUid == "" {
		return fmt.Errorf("file uid required")
	}

	return db.Db().Where("uid = ?", fileUid).Delete(&entity.File{}).Error
}

// DeleteFileUser
func DeleteFileUser(f_u entity.FileUser) error {
	return db.Db().Where("file_id = ? AND user_id = ?", f_u.FileID, f_u.UserID).Delete(&entity.FileUser{}).Error
}