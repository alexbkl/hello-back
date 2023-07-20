package entity

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
	"gorm.io/gorm"
)

type role string

const (
	AdminRole role = "admin"
	UserRole  role = "user"
)

const (
	UserUID = byte('u')
)

type User struct {
	gorm.Model
	UserUID string `gorm:"type:varchar(42);column:user_uid;uniqueIndex"`
	Name    string `gorm:"unique;not null;max:50" json:"name"`
	Role    role   `gorm:"not null;default:user" json:"role"`
}

// TableName returns the entity table name.
func (User) TableName() string {
	return "users"
}

func (m *User) Create() error {
	return db.Db().Create(m).Error
}

// BeforeCreate sets a random UID if needed before inserting a new row to the database.
func (m *User) BeforeCreate(db *gorm.DB) error {

	if rnd.IsUnique(m.UserUID, UserUID) {
		return nil
	}

	m.UserUID = rnd.GenerateUID(UserUID)
	db.Statement.SetColumn("UserUID", m.UserUID)

	return nil
	// return db.Scopes().SetColumn("UserUID", m.UserUID)
}
