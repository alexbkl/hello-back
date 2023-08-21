package query

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
)

// DeleteFileUser
func DeleteFileUser(f_u entity.FileUser) error {
	return db.Db().
		Where("file_id = ? AND user_id = ?", f_u.FileID, f_u.UserID).
		Delete(&entity.FileUser{}).
		Error
}
