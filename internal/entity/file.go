package entity

import (
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
	FileUID   string `gorm:"type:varchar(42);index;" json:"UID"`
	FileName  string `gorm:"type:varchar(1024);uniqueIndex:idx_files_name_root;" json:"name"`
	FileRoot  string `gorm:"type:varchar(255);default:'/';uniqueIndex:idx_files_name_root;" json:"root"`
	FileMime  string `gorm:"type:varchar(64)" json:"mime"`
	FileHash  string `gorm:"type:varchar(128);index" json:"hash"`
	FileSize  int64  `json:"Size"`
	MediaType string `gorm:"type:varchar(16)" json:"MediaType"`
}

// TableName returns the entity table name.
func (File) TableName() string {
	return "files"
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *File) BeforeCreate(db *gorm.DB) error {
	// Set MediaType based on FileName if empty.
	if m.MediaType == "" && m.FileName != "" {
		m.MediaType = media.FromName(m.FileName).String()
	}

	// Return if uid exists.
	if rnd.IsUnique(m.FileUID, FileUID) {
		return nil
	}

	db.Statement.SetColumn("FileUID", rnd.GenerateUID(FileUID))

	return nil
}
