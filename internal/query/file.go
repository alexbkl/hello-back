package query

import (
	"fmt"

	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
)

// FileByUID returns file for the given UID.
func FindFileByUID(uid string) (*entity.File, error) {
	f := entity.File{}

	if uid == "" {
		return &f, fmt.Errorf("file uid required")
	}

	err := db.Db().Where("uid = ?", uid).First(&f).Error

	return &f, err
}

// FilesByRoot return files in a given folder root.
func FindFilesByRoot(root string) (files entity.Files, err error) {
	if err := db.Db().Where("root = ?", root).Find(&files).Error; err != nil {
		return files, err
	}

	return files, err
}

func FindRootFilesByUser(user_id uint) (files entity.Files, err error) {
	if err := db.Db().
		Table("files").
		Joins("LEFT JOIN files_users on files_users.file_id = files.id").
		Where("files.root = '/' AND files_users.permission = 'owner' AND files_users.user_id = ?", user_id).
		Find(&files).Error; err != nil {
		return files, err
	}

	return files, nil
}

// DeleteFileByUID deletes a file by its UID.
func DeleteFileByUID(fileUid string) error {
	if fileUid == "" {
		return fmt.Errorf("file uid required")
	}

	return db.Db().Where("uid = ?", fileUid).Delete(&entity.File{}).Error
}
