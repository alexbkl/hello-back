package entity

import (
	"time"

	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/pkg/rnd"
	"gorm.io/gorm"
)

type role string

type Users []User

const (
	AdminRole role = "admin"
	UserRole  role = "user"
)

const (
	UserUID = byte('u')
)

type User struct {
	ID        uint           `gorm:"primarykey"                   json:"id"`
	UID       string         `gorm:"type:varchar(42);uniqueIndex" json:"uid"`
	Name      string         `gorm:"not null;max:50"              json:"name"`
	Role      role           `gorm:"not null;default:user"        json:"role"`
	Email     Email          `                                    json:"email"`
	Wallet    Wallet         `                                    json:"wallet"`
	Github    Github         `                                    json:"github"`
	CreatedAt time.Time      `                                    json:"created_at"`
	UpdatedAt time.Time      `                                    json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                        json:"deleted_at"`
}

// TableName returns the entity table name.
func (User) TableName() string {
	return "users"
}

func (m *User) Create() error {
	return db.Db().Create(m).Error
}

func (m *User) Save() error {
	return db.Db().Save(m).Error
}

// BeforeCreate sets a random UID if needed before inserting a new row to the database.
func (m *User) BeforeCreate(db *gorm.DB) error {
	if rnd.IsUnique(m.UID, UserUID) {
		return nil
	}

	m.UID = rnd.GenerateUID(UserUID)
	db.Statement.SetColumn("UID", m.UID)

	return nil
	// return db.Scopes().SetColumn("UserUID", m.UserUID)
}

func (m *User) RetrieveNonce(renew bool) (string, error) {
	u := &User{}
	w := &Wallet{}

	// query for find user from wallet address
	if err := db.Db().Model(&u).Preload("Wallet").Where("id IN (?)", db.Db().Table("wallets").Select("user_id").Where("address = ?", m.Wallet.Address)).First(&u).Error; err == nil {

		log.Info("err: ", err)
		w = &u.Wallet
		if renew {
			w.Nonce = rnd.GenerateRandomString(16)
			if err := w.Save(); err != nil {
				return "", err
			}
		}
		return w.Nonce, nil
	} else {
		m.Name = m.Wallet.Address

		if err := m.Create(); err != nil {
			return "", err
		}
	}

	return m.Wallet.Nonce, nil
}
