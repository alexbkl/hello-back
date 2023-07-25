package query

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
)

// RegisteredUsers finds all registered users.
func RegisteredUsers() (result entity.Users) {
	if err := db.Db().Where("id > 0").Find(&result).Error; err != nil {
		log.Errorf("users: %s", err)
	}

	return result
}

func FindUser(find entity.User) *entity.User {
	m := &entity.User{}

	stmt := db.UnscopedDb()

	if find.ID != 0 && find.Name != "" {
		stmt = stmt.Where("id = ? OR user_name = ?", find.ID, find.Name)
	} else if find.ID != 0 {
		stmt = stmt.Where("id = ?", find.ID)
	} else if rnd.IsUID(find.UserUID, entity.UserUID) {
		stmt = stmt.Where("user_uid = ?", find.UserUID)
	} else if find.Name != "" {
		stmt = stmt.Where("user_name = ?", find.Name)
	} else {
		return nil
	}

	// Find matching record.
	if err := stmt.First(m).Error; err != nil {
		return nil
	}

	return m

}
