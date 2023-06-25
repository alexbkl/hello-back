package handlers

import (
	"errors"
	"fmt"
	"meta-go-api/config"
	"meta-go-api/entities"
	"strings"

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
		fmt.Printf("user already has a password")
		//compare password with existing password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
		if err != nil {
			fmt.Printf("Error comparing password: %s", err.Error())
			return c.Status(503).SendString("Passwords don't match")
		} else {
			fmt.Printf("Passwords match")
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
