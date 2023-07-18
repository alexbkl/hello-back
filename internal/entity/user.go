package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type role string

const (
	AdminRole role = "admin"
	UserRole  role = "user"
)

type User struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Name      string    `gorm:"unique;not null;max:50" json:"name"`
	Role      role      `gorm:"type:role;default:user" json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
