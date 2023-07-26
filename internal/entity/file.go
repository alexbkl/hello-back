package entity

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/pkg/media"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
	"gorm.io/gorm"
)

const (
	FileUID = byte('f')
)

// Files represents a file result set.
type Files []File

type File struct {
	gorm.Model
	UID       string `gorm:"type:varchar(42);index;" json:"UID"`
	Name      string `gorm:"type:varchar(1024);" json:"name"`
	Root      string `gorm:"type:varchar(42);default:'/';" json:"root"` // parent folder uid
	Mime      string `gorm:"type:varchar(64)" json:"mime"`
	Size      int64  `json:"Size"`
	MediaType string `gorm:"type:varchar(16)" json:"MediaType"`
}

// TableName returns the entity table name.
func (File) TableName() string {
	return "files"
}

func (m *File) Create() error {
	return db.Db().Create(m).Error
}

func (m *File) Save() error {
	return db.Db().Save(m).Error
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *File) BeforeCreate(db *gorm.DB) error {
	// Set MediaType based on FileName if empty.
	if m.MediaType == "" && m.Name != "" {
		m.MediaType = media.FromName(m.Name).String()
	}

	// Return if uid exists.
	if rnd.IsUnique(m.UID, FileUID) {
		return nil
	}

	db.Statement.SetColumn("UID", rnd.GenerateUID(FileUID))

	return nil
}
