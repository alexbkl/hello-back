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
	FolderUID   string `gorm:"type:varchar(42);index;" json:"UID"`
	FolderTitle string `gorm:"type:varchar(255);" json:"Title"`
	FolderPath  string `gorm:"type:varchar(1024);" json:"Path"`           // folderA/folderB/***
	FolderRoot  string `gorm:"type:varchar(42);default:'/';" json:"Root"` // parent folder uid
}

// TableName returns the entity table name.
func (Folder) TableName() string {
	return "folders"
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *Folder) BeforeCreate(db *gorm.DB) error {
	if rnd.IsUnique(m.FolderUID, 'd') {
		return nil
	}

	db.Statement.SetColumn("FolderUID", rnd.GenerateUID(FolderUID))

	return nil
}
