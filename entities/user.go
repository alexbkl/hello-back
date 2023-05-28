package entities

import (
	"gorm.io/gorm"
)


type User struct {  
	gorm.Model  
	Address string `gorm:"unique;not null" json:"address"`
	Nonce   string `json:"nonce"`

	Files []File `gorm:"foreignKey:UserAddress;references:Address"`  
  }  
  

