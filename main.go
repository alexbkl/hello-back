package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"meta-go-api/config"
	"meta-go-api/handlers"
)

//say helloworld
func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	var err error
	if config.Database, err = config.Connect(); err != nil {
		panic(err)
	}

	app.Post("/register", handlers.RegisterHandler)

	//r.Get("/users/{address:^0x[a-fA-F0-9]{40}$}/nonce", handlers.UserNonceHandler)

	//In Fiber using regex:
	//EXAMPLE:
	app.Get("/users/:address/nonce", handlers.UserNonceHandler)

	app.Post("/signin", handlers.SigninHandler)

	//create a /api group with cors, handlers.AuthMiddleware and handlers.WelcomeHandler
	api := app.Group("/api")

	api.Use(handlers.AuthMiddleware)

	api.Get("/welcome", handlers.WelcomeHandler)

	//doggos
	app.Get("/dogs", handlers.GetDogs)
	app.Get("/gods/:id", handlers.GetDog)
	app.Post("/dogs", handlers.CreateDog)
	app.Put("/dogs/:id", handlers.UpdateDog)
	app.Delete("/dogs/:id", handlers.DeleteDog)

	fmt.Println("Server started")

	log.Fatal(app.Listen(":8001"))

}
