package entities

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	EncryptedMetadata string `json:"encryptedMetadata" gorm:"not null;max:1000"`
	UserAddress string `json:"userAddress" gorm:"not null;max:255"`
	CIDOfEncryptedBuffer string `json:"cidOfEncryptedBuffer" gorm:"not null;max:255"`
	CIDEncryptedOriginalStr string `json:"cidEncryptedOriginalStr" gorm:"not null;max:255"`
	IV string `json:"iv" gorm:"not null;max:255"`
	BytesLength int `json:"bytesLength" gorm:"not null;max:255"`
}