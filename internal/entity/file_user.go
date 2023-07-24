package entity

import (
	"time"

	"github.com/Hello-Storage/hello-back/internal/db"
)

type permission string

const (
	OwnerPermission  role = "owner"
	SharedPermission role = "shared"
)

// FileUser represents a one-to-many relation between File and User.

type FileUser struct {
	FileID     uint       `gorm:"primary_key;auto_increment:false"`
	UserID     uint       `gorm:"primary_key;auto_increment:false"`
	Permission permission `gorm:"not null;" json:"permission"`
	File       *File
	User       *User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName returns the entity table name.
func (FileUser) TableName() string {
	return "files_users"
}

func (m *FileUser) Create() error {
	return db.Db().Create(m).Error
}
