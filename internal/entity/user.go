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
	Detail    UserDetail     `                                    json:"detail"`
	CreatedAt time.Time      `                                    json:"created_at"`
	UpdatedAt time.Time      `                                    json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                        json:"deleted_at"`
}

// TableName returns the entity table name.
func (User) TableName() string {
	return "users"
}

func (user *User) Create() error {
	
	return db.Db().Create(user).Error
}

func (user *User) Save() error {
	return db.Db().Save(user).Error
}

// BeforeCreate sets a random UID if needed before inserting a new row to the database.
func (user *User) BeforeCreate(db *gorm.DB) error {
	if rnd.IsUnique(user.UID, UserUID) {
		return nil
	}

	user.UID = rnd.GenerateUID(UserUID)
	db.Statement.SetColumn("UID", user.UID)

	return nil
	// return db.Scopes().SetColumn("UserUID", m.UserUID)
}

func (user *User) RetrieveNonce(renew bool) (string, error) {
	u := &User{}
	w := &Wallet{}

	// query for find user from wallet address
	subquery := db.Db().Table("wallets").Select("user_id").Where("address = ?", user.Wallet.Address)
	if err := db.Db().Model(&u).Preload("Wallet").Where("id IN (?)", subquery).First(&u).Error; err == nil {
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
		user.Name = user.Wallet.Address

		if err := user.Create(); err != nil {
			return "", err
		}

		// initialize user detail
		user_detail := UserDetail{
			StorageUsed: 0,
			UserID:      user.ID,
		}

		if err := user_detail.Create(); err != nil {
			return "", err
		}
	}

	return user.Wallet.Nonce, nil
}
