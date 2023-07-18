package entities

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	EncryptedMetadata       string `json:"encryptedMetadata" gorm:"not null;max:10000"`
	UserAddress             string `json:"userAddress" gorm:"not null;max:255"`
	CIDOfEncryptedBuffer    string `json:"cidOfEncryptedBuffer" gorm:"not null;max:2550"`
	CIDEncryptedOriginalStr string `json:"cidEncryptedOriginalStr" gorm:"not null;max:2550"`
	IV                      string `json:"iv" gorm:"not null;max:255"`
	BytesLength             int    `json:"bytesLength" gorm:"not null;max:255"`
}

type PublishedFile struct {
	gorm.Model
	Metadata             string `json:"metadata" gorm:"not null;max:10000"`
	CIDOfEncryptedBuffer string `json:"cidOfEncryptedBuffer" gorm:"not null;max:2550"`
	CIDOriginalStr       string `json:"cidOriginalStr" gorm:"not null;max:2550"`
	UserAddress          string `json:"userAddress" gorm:"not null;max:255"`
	Hash                 string `json:"hash" gorm:"unique;not null;max:255"`
}
