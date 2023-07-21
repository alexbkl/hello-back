package entity

import (
	"github.com/Hello-Storage/hello-back/pkg/rnd"
	"gorm.io/gorm"
)

const (
	FolderUID = byte('d')
)

type Folders []Folder

type Folder struct {
	gorm.Model
	UID   string `gorm:"type:varchar(42);index;" json:"UID"`
	Title string `gorm:"type:varchar(255);" json:"Title"`
	Path  string `gorm:"type:varchar(1024);default:'/';" json:"Path"` // folderA/folderB/***
	Root  string `gorm:"type:varchar(42);default:'/';" json:"Root"`   // parent folder uid
}

// TableName returns the entity table name.
func (Folder) TableName() string {
	return "folders"
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *Folder) BeforeCreate(db *gorm.DB) error {
	if rnd.IsUnique(m.UID, 'd') {
		return nil
	}

	db.Statement.SetColumn("UID", rnd.GenerateUID(FolderUID))

	return nil
}
