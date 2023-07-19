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
	BytesLength             uint64 `json:"bytesLength" gorm:"not null;"`
}

type PublishedFile struct {
	gorm.Model
	Metadata             string `json:"metadata" gorm:"not null;max:10000"`
	CIDOfEncryptedBuffer string `json:"cidOfEncryptedBuffer" gorm:"not null;max:2550"`
	CIDOriginalStr       string `json:"cidOriginalStr" gorm:"not null;max:2550"`
	UserAddress          string `json:"userAddress" gorm:"not null;max:255"`
	Hash                 string `json:"hash" gorm:"unique;not null;max:255"`
	File                 File   `gorm:"foreignKey:FileID;references:ID"`
	FileID               uint   `json:"fileID" gorm:"not null"`
}

type FileSharedState struct {
	gorm.Model
	UserAddress     string        `json:"userAddress" gorm:"not null;max:255"`
	PublishedFile   PublishedFile `gorm:"foreignKey:PublishedFileID;references:ID;OnDelete:SET NULL"`
	PublishedFileID *uint         `json:"publishedFileID" gorm:"uniqueIndex"`
	File            File          `gorm:"foreignKey:FileID;references:ID"`
	FileID          uint          `json:"fileID" gorm:"not null"`
}

type OneTimeFile struct {
	gorm.Model
	Visited         bool          `json:"visited" gorm:"not null;default:false"`
	PublishedFileID uint          `json:"publishedFileID" gorm:"not null"`
	PublishedFile   PublishedFile `gorm:"foreignKey:PublishedFileID;references:ID"`
}

/*
	OneTime bool `json:"oneTime"`
	AddressRestricted bool `json:"addressRestricted"`
	PasswordProtected bool `json:"passwordProtected"`
	TemporaryLink bool `json:"temporaryLink"`
	Subscription bool `json:"subscription"`
	SubscriptionPrice uint `json:"subscriptionPrice"`
	SubscriptionPeriod uint `json:"subscriptionPeriod"`
*/
