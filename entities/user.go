package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
}