package handlers

import (
	"errors"

	"github.com/Hello-Storage/hello-back/config"
	"github.com/Hello-Storage/hello-back/entities"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidEmail = errors.New("email is invalid")
)

type emailPayload struct {
	Email string `json:"email"`
}

func (p *emailPayload) Validate() error {
	if p.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

func SubmitEmail(c *fiber.Ctx) error {
	var p = new(emailPayload)

	if err := c.BodyParser(p); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	if err := p.Validate(); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	var email entities.Email

	//where email.Email = p.Email
	result := config.Database.Find(&email, "email = ?", p.Email)

	if result.RowsAffected != 0 {
		return c.Status(409).SendString("email already exists")
	}

	e := entities.Email{
		Email: p.Email,
	}

	config.Database.Create(&e)

	return c.Status(201).JSON(e)
}
