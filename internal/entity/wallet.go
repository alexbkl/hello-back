package entity

import (
	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
	"gorm.io/gorm"
)

type Wallet struct {
	ID      uint   `gorm:"primarykey"`
	Address string `gorm:"type:varchar(50);not null;uniqueIndex" json:"address"`
	Type    string `gorm:"type:varchar(30);not null;default:eth" json:"type"`
	Nonce   string `gorm:"type:varchar(16);not null" json:"nonce"`
	UserID  uint
}

// TableName returns the entity table name.
func (Wallet) TableName() string {
	return "wallets"
}

func (m *Wallet) Create() error {
	return db.Db().Create(m).Error
}

// BeforeCreate creates a random UID if needed before inserting a new row to the database.
func (m *Wallet) BeforeCreate(db *gorm.DB) error {
	m.Nonce = rnd.GenerateRandomString(16)
	db.Statement.SetColumn("nonce", m.Nonce)
	return nil
}

func (m *Wallet) Save() error {
	return db.Db().Save(m).Error
}
