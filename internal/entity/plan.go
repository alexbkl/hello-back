package entity

import "github.com/Hello-Storage/hello-back/internal/db"

type Plan struct {
	ID               uint   `gorm:"primarykey"              json:"id"`
	Name             string `gorm:"type:varchar(50);unique" json:"name"`
	StorageAvailable uint   `                               json:"storage_available"` // bytes format

}

// TableName returns the entity table name.
func (Plan) TableName() string {
	return "plans"
}

func (m *Plan) Create() error {
	return db.Db().Create(m).Error
}

func (m *Plan) Save() error {
	return db.Db().Save(m).Error
}
