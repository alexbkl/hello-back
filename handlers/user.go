package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Hello-Storage/hello-back/config"
	"github.com/Hello-Storage/hello-back/entities"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("password is invalid")
)

type passwordPayload struct {
	Password string `json:"password"`
}

func (p *passwordPayload) Validate() error {
	//if pass is empty, does not contain at least 8 characters, does not contain at least 1 uppercase letter, does not contain at least 1 lowercase letter, does not contain at least 1 number, does not contain at least 1 special character
	if p.Password == "" {
		return errors.New("password is required")
	} else if len(p.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	} else if !strings.ContainsAny(p.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.New("password must contain at least 1 uppercase letter")
	} else if !strings.ContainsAny(p.Password, "abcdefghijklmnopqrstuvwxyz") {
		return errors.New("password must contain at least 1 lowercase letter")
	} else if !strings.ContainsAny(p.Password, "0123456789") {
		return errors.New("password must contain at least 1 number")
	} else if !strings.ContainsAny(p.Password, "!@#$%^&*()_+-=,./<>?;:'\"[]{}\\|`~") {
		return errors.New("password must contain at least 1 special character")
	}

	return nil
}

func SubmitPasswordHandler(c *fiber.Ctx) error {
	var p = new(passwordPayload)

	if err := c.BodyParser(p); err != nil {
		fmt.Printf("Parsing error: %s", err.Error())
		return c.Status(503).SendString("Parsing error: " + err.Error())
	}

	if err := p.Validate(); err != nil {
		fmt.Printf("Validation error: %s", err.Error())
		return c.Status(503).SendString("Validation error: " + err.Error())
	}

	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	//hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)

	//check if user already has a password
	if user.Password != "" {
		//fmt.Printf("user already has a password")
		//compare password with existing password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
		if err != nil {
			fmt.Printf("Error comparing password: %s", err.Error())
			return c.Status(503).SendString("Passwords don't match")
		} else {
			//fmt.Printf("Passwords match")
			return c.Status(200).SendString("Signed in successfully")
		}
	}

	if err != nil {
		fmt.Printf("Error hashing password: %s", err.Error())
		return c.Status(503).SendString(err.Error())
	}

	//update user password
	user.Password = string(hashedPassword)

	//update user in database
	config.Database.Save(&user)

	return c.Status(200).SendString("Registered successfully")

}

func GetDatacapHandler(c *fiber.Ctx) error {
	var user entities.User

	address := c.Params("address")

	result := config.Database.Where("address = ?", address).First(&user)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	} else if result.Error != nil {
		return c.Status(503).SendString(result.Error.Error())
	}

	return c.Status(200).JSON(user.DataCap)
}

func GetUsedStorageHandler(c *fiber.Ctx) error {
	var user entities.User

	address := c.Params("address")

	result := config.Database.Where("address = ?", address).First(&user)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	} else if result.Error != nil {
		return c.Status(503).SendString(result.Error.Error())
	}

	return c.Status(200).JSON(user.UsedStorage)
}

func GetTotalUploadedFilesHandler(c *fiber.Ctx) error {
	var user entities.User

	address := c.Params("address")

	result := config.Database.Where("address = ?", address).First(&user)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	} else if result.Error != nil {
		return c.Status(503).SendString(result.Error.Error())
	}

	return c.Status(200).JSON(user.TotalUploadedFiles)
}

/*
type File struct {
	gorm.Model
	EncryptedMetadata string `json:"encryptedMetadata" gorm:"unique;not null;max:1000"`
	UserAddress string `json:"userAddress" gorm:"unique;not null;max:255"`
	CIDOfEncryptedBuffer string `json:"cidOfEncryptedBuffer" gorm:"unique;not null;max:255"`
	CIDEncryptedOriginalStr string `json:"cidEncryptedOriginalStr" gorm:"unique;not null;max:255"`
	IV string `json:"iv" gorm:"unique;not null;max:255"`
	BytesLength int `json:"bytesLength" gorm:"not null;max:255"`
}
*/

/*
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
*/

func GetUploadedFilesCountHandler(c *fiber.Ctx) error {
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	var count int64

	result := config.Database.Model(&entities.File{}).Where("user_address = ?", user.Address).Count(&count)

	if result.Error != nil {
		fmt.Printf("Error getting number of uploaded files: %s", result.Error.Error())
		return c.Status(503).SendString(result.Error.Error())
	}

	if count == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(count)
}
