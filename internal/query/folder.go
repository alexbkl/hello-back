package query

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
)

func FindFolder(find entity.Folder) *entity.Folder {
	m := &entity.Folder{}

	stmt := db.Db()

	if find.ID != 0 && find.Title != "" {
		stmt = stmt.Where("id = ? OR title = ?", find.ID, find.Title)
	} else if find.ID != 0 {
		stmt = stmt.Where("id = ?", find.ID)
	} else if rnd.IsUID(find.UID, entity.FolderUID) {
		stmt = stmt.Where("uid = ?", find.UID)
	} else if find.Title != "" {
		stmt = stmt.Where("title = ?", find.Title)
	} else {
		return nil
	}

	// Find matching record.
	if err := stmt.First(m).Error; err != nil {
		return nil
	}

	return m

}

// FoldersByRoot returns folders in a given directory.
func FoldersByRoot(root string) (folders entity.Folders, err error) {
	if err := db.Db().Where("root = ?", root).Find(&folders).Error; err != nil {
		return folders, err
	}

	return folders, nil
}

func FindRootFoldersByUser(user_id uint) (folders entity.Folders, err error) {
	if err := db.Db().
		Table("folders").
		Joins("LEFT JOIN folders_users on folders_users.folder_id = folders.id").
		Where("folders.root = '/' AND folders_users.permission = 'owner' AND folders_users.user_id = ?", user_id).
		Find(&folders).Error; err != nil {
		return folders, err
	}

	return folders, nil
}

func FindFolderByTitleAndRoot(title, root string) *entity.Folder {
	m := &entity.Folder{}

	stmt := db.Db()
	stmt = stmt.Where("title = ? AND root = ?", title, root)

	// Find matching record.
	if err := stmt.First(m).Error; err != nil {
		return nil
	}

	return m
}

func FindFolderPathByRoot(root string) entity.Folders {
	if root == "/" {
		return entity.Folders{}
	}

	m := FindFolder(entity.Folder{UID: root})

	return append(FindFolderPathByRoot(m.Root), *m)
}
