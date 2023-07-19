package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Hello-Storage/hello-back/config"
	"github.com/Hello-Storage/hello-back/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PublishPayload struct {
	CidOriginalStr       string `json:"cidOriginalStr"` //this is the cid of the original buffer
	CidOfEncryptedBuffer string `json:"cidOfEncryptedBuffer"`
	Metadata             string `json:"metadata"`
	FileID               uint   `json:"fileID"`
}

var (
	ErrInvalidCidOriginalBuffer    = errors.New("cidOriginalBuffer is invalid")
	ErrInvalidCidOfEncryptedBuffer = errors.New("cidOfEncryptedBuffer is invalid")
	ErrInvalidMetadata             = errors.New("metadata is invalid")
	ErrInvalidFileID               = errors.New("fileID is invalid")
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
	if p.FileID == 0 {
		return ErrInvalidFileID
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

func getSharedFile(fileId string, userAddress string) (*entities.FileSharedState, error) {
	var fileSharedState entities.FileSharedState
	//find FileSharedState where file's id is equal to the id of the file in the request
	//prevent log
	result := config.Database.Preload("PublishedFile").Where("file_id = ? AND user_address = ?", fileId, userAddress).First(&fileSharedState)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//No error, just no result
			return nil, nil
		} else {
			fmt.Printf("Error getting fileSharedState: %s", result.Error.Error())
			return nil, result.Error
		}
	}

	return &fileSharedState, nil
}

func PublishFileHandler(c *fiber.Ctx) error {
	//store  details in database and generate hash based on content, then return link (hash) to user
	var p PublishPayload
	var user entities.User
	var publishedFileExists entities.PublishedFile

	//get user from context
	user = c.Locals("user").(entities.User)

	if err := c.BodyParser(&p); err != nil {
		fmt.Println("Error parsing body: ", err)
		return c.Status(400).SendString("Parsing error: " + err.Error())
	}

	if err := p.Validate(); err != nil {
		fmt.Println("Error validating body: ", err)
		return c.Status(400).SendString("Validation error: " + err.Error())
	}
	//check if file already exists
	result := config.Database.Where("c_id_of_encrypted_buffer = ? and user_address = ?", p.CidOfEncryptedBuffer, user.Address).Find(&publishedFileExists)
	if result.RowsAffected > 0 {
		fmt.Println("File already published")
		return c.Status(400).SendString("File already published")
	}
	//create hash based on content

	hash := createHash(p)

	//store details in database and generate hash based on content, then return link (hash) to user
	publishedFile := entities.PublishedFile{
		Metadata:             p.Metadata,
		CIDOriginalStr:       p.CidOriginalStr,
		CIDOfEncryptedBuffer: p.CidOfEncryptedBuffer,
		UserAddress:          user.Address,
		FileID:               p.FileID,
		Hash:                 hash,
	}

	//create published file
	config.Database.Create(&publishedFile)
	
	//return link to user
	return c.Status(200).JSON(fileSharedState)

}

func UnpublishFileHandler(c *fiber.Ctx) error {
	//delete file from database
	var user entities.User
	var publishedFile entities.PublishedFile
	var fileSharedState entities.FileSharedState

	//get user from context
	user = c.Locals("user").(entities.User)

	fileID := c.Params("fileId")

	//check if file id is valid
	if fileID == "" {
		fmt.Println("File ID is invalid")
		return c.Status(400).SendString(fmt.Sprintf("File ID %s is invalid", fileID))
	}

	//check that file exists and belongs to user's address
	fileIDInt, err := strconv.Atoi(fileID)
	if err != nil {
		fmt.Println("File ID is invalid")
		return c.Status(400).SendString(fmt.Sprintf("File ID %s is invalid", fileID))
	}
	ownership := config.Database.Where("file_id = ? AND user_address = ?", fileIDInt, user.Address).Find(&publishedFile)

	if ownership.RowsAffected == 0 {
		fmt.Printf("File with ID %s does not exist or does not belong to user with address %s", fileID, user.Address)
		if ownership.Error != nil {
			fmt.Printf("Error checking ownership: %s", ownership.Error.Error())
		}
		return c.Status(404).SendString(fmt.Sprintf("File with ID %s does not exist or does not belong to user", fileID))
	}

	//reset file shared state
	result := config.Database.Where("file_id = ? AND user_address = ?", fileID, user.Address).Find(&fileSharedState)

	if result.RowsAffected != 0 {
		fileSharedState.PublishedFileID = nil
		config.Database.Save(&fileSharedState)

	}
	//delete file from database
	result = config.Database.Unscoped().Delete(&publishedFile)

	if result.Error != nil {
		fmt.Printf("Error deleting file: %s", result.Error.Error())
		return c.Status(503).SendString(result.Error.Error())
	}

	//reload fileSharedState without PublishedFile
	fileSharedStatePtr, err := getSharedFile(fileID, user.Address)
	if err != nil {
		fmt.Printf("Error reloading file shared state: %s", err.Error())
		return c.Status(503).SendString("Error reloading file shared state: " + err.Error())
	}

	return c.Status(200).JSON(*fileSharedStatePtr)

}

func OneTimeShareHandler(c *fiber.Ctx) error {
	//get user from context
	var p PublishPayload
	var publishedFileExists entities.PublishedFile

	user := c.Locals("user").(entities.User)

	if err := c.BodyParser(&p); err != nil {
		fmt.Println("Error parsing body: ", err)
		return c.Status(400).SendString("Error parsing body: " + err.Error())
	}

	if err := p.Validate(); err != nil {
		fmt.Println("Error validating body: ", err)
		return c.Status(400).SendString("Validation error: " + err.Error())
	}

	//check if file already exists
	result := config.Database.Where("c_id_of_encrypted_buffer = ? and user_address = ?", p.CidOfEncryptedBuffer, user.Address).Find(&publishedFileExists)
	if result.RowsAffected != 0 {
		fmt.Println("File already published")
		return c.Status(400).SendString("File already published")
	}
	//check if a one-time file with this file ID already exists
	var existingOneTimeFile entities.OneTimeFile
	if err := config.Database.Where("file_id = ?", p.FileID).First(&existingOneTimeFile).Error; err == nil {
		//if a one-time file already exists, return an error
		return c.Status(400).SendString(fmt.Sprintf("One-time share for file ID %d already exists", p.FileID))
	}

	//check that the file exists and belongs to the user's address
	var file entities.File
	if err := config.Database.Where("id = ? AND user_address = ?", p.FileID, user.Address).First(&file).Error; err != nil {
		return c.Status(404).SendString(fmt.Sprintf("File ID %s not found or does not belong to user", p.CidOfEncryptedBuffer))
	}

	hash := createHash(p)

	publishedFile := entities.PublishedFile{
		Metadata:             p.Metadata,
		CIDOriginalStr:       p.CidOriginalStr,
		CIDOfEncryptedBuffer: p.CidOfEncryptedBuffer,
		UserAddress:          user.Address,
		FileID:               p.FileID,
		Hash:                 hash,
	}

	//create published file
	if err := config.Database.Create(&publishedFile).Error; err != nil {
		fmt.Printf("Error creating published file: %s", err.Error())
		return c.Status(503).SendString("Error creating published file: " + err.Error())
	}

	oneTimeFile := entities.OneTimeFile{
		Visited:         false,
		PublishedFileID: publishedFile.ID,
	}

	//create one-time file
	if err := config.Database.Create(&oneTimeFile).Error; err != nil {
		fmt.Printf("Error creating one-time file: %s", err.Error())
		return c.Status(503).SendString("Error creating one-time file: " + err.Error())
	}

	//Get FileSharedState
	fileSharedState, err := getSharedFile(strconv.Itoa(int(p.FileID)), user.Address)
	if err != nil || fileSharedState == nil {
		// create file shared state
		fileSharedState = &entities.FileSharedState{
			UserAddress:     user.Address,
			PublishedFileID: &publishedFile.ID,
			FileID:          p.FileID,
		}
		if err := config.Database.Create(fileSharedState).Error; err != nil {
			fmt.Println("Error creating file shared state: ", err)
			return c.Status(400).SendString("Error creating file shared state: " + err.Error())
		}
	} else {
		//update file shared state
		fileSharedState.PublishedFileID = &publishedFile.ID
		if err := config.Database.Save(fileSharedState).Error; err != nil {
			fmt.Println("Error updating file shared state: ", err)
			return c.Status(400).SendString("Error updating file shared state: " + err.Error())
		}
	}

	//reload fileSharedState with PublishedFile
	fileSharedState, err = getSharedFile(strconv.Itoa(int(p.FileID)), user.Address)
	if err != nil || fileSharedState == nil {
		fmt.Println("Error reloading file shared state: ", err)
		return c.Status(400).SendString("Error reloading file shared state: " + err.Error())
	}

	return c.Status(200).JSON(fileSharedState)

}
/*
// to check if a file is visited and delete it if it is
func (o *entities.OneTimeFile) CheckAndDelete() error {
	if o.Visited {
		//if the file has been visited, delete it
		if err := config.Database.Delete(o).Error; err != nil {
			return err
		}
	}
	return nil
}

// to mark file as visited
func (o *entities.OneTimeFile) MarkAsVisited() error {
	o.Visited = true
	if err := config.Database.Save(o).Error; err != nil {
		return err
	}
	return nil
}
*/
func GetSharedFileStateHandler(c *fiber.Ctx) error {
	var user entities.User
	var file entities.File

	//var publishedFile entities.PublishedFile
	fileID := c.Params("fileId")
	//get user from context
	user = c.Locals("user").(entities.User)

	//check if file id is valid
	if fileID == "" {
		return c.Status(400).SendString("File ID is invalid")
	}

	//check that file exists and belongs to user's address
	if err := config.Database.Where("id = ? AND user_address = ?", fileID, user.Address).First(&file).Error; err != nil {
		return c.Status(404).SendString("File does not exist or does not belong to user")
	}

	//get FileSharedState
	fileSharedState, err := getSharedFile(fileID, user.Address)
	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	//return FileSharedState
	return c.Status(200).JSON(fileSharedState)

}

func GetPublishedFileHandler(c *fiber.Ctx) error {
	//get the hash from the URL parameters
	hash := c.Params("hash")

	// if hash is empty, return an error
	if hash == "" {
		return c.Status(400).SendString("Hash is invalid")
	}

	//query the database for the published file with the matching hash
	var publishedFile entities.PublishedFile
	result := config.Database.Where("hash = ?", hash).First(&publishedFile)

	//if the file was not found, return an error
	if result.Error != nil {
		fmt.Printf("Error retrieving file: %s", result.Error.Error())
		return c.Status(404).SendString("File not found")
	}

	//if the file was found, return the metadata
	return c.Status(200).JSON(publishedFile)
}

