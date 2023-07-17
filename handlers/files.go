package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Hello-Storage/hello-back/config"
	"github.com/Hello-Storage/hello-back/entities"

	"github.com/gofiber/fiber/v2"
)

type PublishPayload struct {
	CidOriginalStr       string `json:"cidOriginalStr"` //this is the cid of the original buffer
	CidOfEncryptedBuffer string `json:"cidOfEncryptedBuffer"`
	Metadata             string `json:"metadata"`
}

var (
	ErrInvalidCidOriginalBuffer    = errors.New("cidOriginalBuffer is invalid")
	ErrInvalidCidOfEncryptedBuffer = errors.New("cidOfEncryptedBuffer is invalid")
	ErrInvalidMetadata             = errors.New("metadata is invalid")
)

func (p *PublishPayload) Validate() error {
	if p.CidOfEncryptedBuffer == "" {
		return ErrInvalidCidOfEncryptedBuffer
	}
	if p.Metadata == "" {
		return ErrInvalidMetadata
	}
	if p.CidOriginalStr == "" {
		return ErrInvalidCidOriginalBuffer
	}

	return nil
}

func createHash(p PublishPayload) string {
	h := sha256.New()
	h.Write([]byte(p.CidOriginalStr))
	h.Write([]byte(p.CidOfEncryptedBuffer))
	h.Write([]byte(p.Metadata))
	return hex.EncodeToString(h.Sum(nil))
}

func PublishFileHandler(c *fiber.Ctx) error {
	//store  details in database and generate hash based on content, then return link (hash) to user
	var p PublishPayload
	var user entities.User

	//get user from context
	user = c.Locals("user").(entities.User)

	if err := c.BodyParser(&p); err != nil {
		return c.Status(503).SendString("Parsing error: " + err.Error())
	}

	if err := p.Validate(); err != nil {
		return c.Status(503).SendString("Validation error: " + err.Error())
	}

	hash := createHash(p)

	fmt.Println("hash: ", hash)
	fmt.Println("cidOriginalStr: ", p.CidOriginalStr)
	fmt.Println("cidOfEncryptedBuffer: ", p.CidOfEncryptedBuffer)
	fmt.Println("metadata: ", p.Metadata)

	//store details in database and generate hash based on content, then return link (hash) to user
	publishedFile := entities.PublishedFile{
		Metadata:             p.Metadata,
		CIDOriginalStr:       p.CidOriginalStr,
		CIDOfEncryptedBuffer: p.CidOfEncryptedBuffer,
		UserAddress:          user.Address,
		Hash:                 hash,
	}

	//create published file
	config.Database.Create(&publishedFile)

	//return link to user
	return c.Status(200).SendString("File published successfully")

}

func GetSharedFileStatesHandler(c *fiber.Ctx) error {
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
