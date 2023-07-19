package entity

import "gorm.io/gorm"

type Folders []Folder

type Folder struct {
	gorm.Model
	FolderUID   string `gorm:"type:varchar(42);index;" json:"UID"`
	Path        string `gorm:"type:varchar(1024);uniqueIndex:idx_folders_path_root;" json:"Path"`
	Root        string `gorm:"type:varchar(16);default:'';uniqueIndex:idx_folders_path_root;" json:"Root"`
	FolderTitle string `gorm:"type:varchar(255);" json:"Title"`
}

// TableName returns the entity table name.
func (Folder) TableName() string {
	return "folders"
}
