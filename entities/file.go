package entities

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	EncryptedMetadata string `json:"encryptedMetadata"`
	UserAddress string `json:"userAddress"`
	CIDOfEncryptedBuffer string `json:"cidOfEncryptedBuffer"`
	CIDEncryptedOriginalStr string `json:"cidEncryptedOriginalStr"`
	IV string `json:"iv"`
}