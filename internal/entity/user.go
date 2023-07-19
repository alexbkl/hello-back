package entity

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type role string

const (
	AdminRole role = "admin"
	UserRole  role = "user"
)

type User struct {
	gorm.Model
	UUID uuid.UUID `gorm:"type:uuid;column:user_uuid;index;default:uuid_generate_v4()"`
	Name string    `gorm:"unique;not null;max:50" json:"name"`
	Role role      `gorm:"not null;default:user" json:"role"`
}

// TableName returns the entity table name.
func (User) TableName() string {
	return "users"
}

func (m *User) Create() error {
	return db.Db().Create(m).Error
}
