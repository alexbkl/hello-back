package entity

import "github.com/Hello-Storage/hello-back/internal/db"

type Github struct {
	ID       uint   `gorm:"primarykey"                   json:"id"`
	GithubID uint   `gorm:"uniqueIndex;column:github_id" json:"github_id"`
	Name     string `gorm:"type:varchar(50);not null"    json:"name"`
	Avatar   string `gorm:"type:varchar(200)"            json:"avatar"`
	UserID   uint
}

// TableName returns the entity table name.
func (Github) TableName() string {
	return "githubs"
}

func (m *Github) Create() error {
	return db.Db().Create(m).Error
}

func (m *Github) Save() error {
	return db.Db().Save(m).Error
}
