package entity

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
	"gorm.io/gorm"
)

const (
	FolderUID = byte('d')
)

type Folders []Folder

type Folder struct {
	gorm.Model
	UID   string `gorm:"type:varchar(42);index;" json:"uid"`
	Title string `gorm:"type:varchar(255);" json:"title"`
	Path  string `gorm:"type:varchar(1024);default:'/';" json:"path"` // folderA/folderB/***
	Root  string `gorm:"type:varchar(42);default:'/';" json:"root"`   // parent folder uid
}

// TableName returns the entity table name.
func (Folder) TableName() string {
	return "folders"
}

func (m *Folder) Create() error {
	return db.Db().Create(m).Error
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *Folder) BeforeCreate(db *gorm.DB) error {
	if rnd.IsUnique(m.UID, 'd') {
		return nil
	}
	m.UID = rnd.GenerateUID(FolderUID)
	db.Statement.SetColumn("UID", m.UID)

	return nil
}
