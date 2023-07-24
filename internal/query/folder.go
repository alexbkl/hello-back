package query

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
)

// FoldersByRoot returns folders in a given directory.
func FoldersByRoot(root string) (folders entity.Folders, err error) {
	if err := db.Db().Where("root = ?", root).Find(&folders).Error; err != nil {
		return folders, err
	}

	return folders, nil
}
