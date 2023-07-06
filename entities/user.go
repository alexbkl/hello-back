package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Address     string `gorm:"unique;not null;max:43" json:"address"`
	Nonce       string `json:"nonce"`
	DataCap     int64  `json:"dataCap"`
	UsedStorage int64  `json:"usedStorage"`
	TotalUploadedFiles int64 `json:"totalUploadedFiles"`
	Files       []File `gorm:"foreignKey:UserAddress;references:Address"`
	Password    string `json:"password"`
}

type Email struct {
	gorm.Model
	Email string `gorm:"unique;not null" json:"email"`
}
