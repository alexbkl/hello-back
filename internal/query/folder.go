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

func FindFolderByTitleWithRoot(title, root string) *entity.Folder {
	m := &entity.Folder{}

	stmt := db.UnscopedDb()
	stmt = stmt.Where("title = ? AND root = ?", title, root)

	// Find matching record.
	if err := stmt.First(m).Error; err != nil {
		return nil
	}

	return m
}
