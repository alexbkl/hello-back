package entities

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	FileName string `json:"filename"`
	UserAddress string `json:"userAddress"`
	CID string `json:"cid"`
}