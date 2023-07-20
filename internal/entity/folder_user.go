package entity

import "time"

// FolderUser represents a one-to-many relation between File and User.

type FolderUser struct {
	FolderID   uint       `gorm:"primary_key;auto_increment:false"`
	UserID     uint       `gorm:"primary_key;auto_increment:false"`
	Permission permission `gorm:"not null;" json:"permission"`
	Folder     *Folder
	User       *User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName returns the entity table name.
func (FolderUser) TableName() string {
	return "folders_users"
}
