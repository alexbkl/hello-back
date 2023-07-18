package main

import (
	"fmt"
	"log"
	"meta-go-api/config"
	"meta-go-api/environment"
	"meta-go-api/handlers"
	"meta-go-api/s3client"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

/*
import (

	"meta-go-api/environment"
	"meta-go-api/s3client"

)
*/
func main() {

	//set environment variables
	environment.SetEnv()

	//init s3client credentials and connection
	s3client.Init()

	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Disposition, Origin, Content-Type, Accept, Authorization, Content-Length, Original-Filename",
		AllowCredentials: true,
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

	api.Post("/personalSignature", handlers.SubmitPersonalSignatureHandler)

	api.Get("/authenticate", handlers.AuthenticationHandler)

	//post method to upload file to s3 and save the file name to database
	api.Post("/file/upload", handlers.UploadHandler)
	//delete method to delete file from s3 and database
	api.Delete("/file/:fileId", handlers.DeleteFileHandler)
	api.Get("/file/:cid", handlers.DownloadFileHandler)

	//get method to get all files from database
	api.Get("/files", handlers.GetFilesHandler)

	//get method to get datacap of a certain address (/datacap/:address)
	api.Get("/storage/datacap/get/:address", handlers.GetDatacapHandler)
	//get method for used storage
	api.Get("/storage/used/get/:address", handlers.GetUsedStorageHandler)
	//get method for uploaded files count
	api.Get("/storage/uploaded/count/get/:address", handlers.GetUploadedFilesCountHandler)

	//get shared file states for frontend
	api.Get("/v0/file/share/states/:fileId", handlers.GetSharedFileStateHandler)
	


	//public
	//publish file (post method that returns personalized sharable URL as long as it exists in database)
	api.Post(("/v0/file/share/publish"), handlers.PublishFileHandler)
	api.Post(("/v0/file/share/one-time/"), handlers.OneTimeShareHandler)
	api.Delete(("/v0/file/share/unpublish/:fileId"), handlers.UnpublishFileHandler)
	//get method to get metadata of public file based on hash from database
	api.Get("/v0/file/share/file/:hash", handlers.GetPublishedFileHandler)


	//doggos
	app.Get("/dogs", handlers.GetDogs)
	app.Get("/gods/:id", handlers.GetDog)
	app.Post("/dogs", handlers.CreateDog)
	app.Put("/dogs/:id", handlers.UpdateDog)
	app.Delete("/dogs/:id", handlers.DeleteDog)
	app.Post("email/submit", handlers.SubmitEmail)
	fmt.Println("Server started")

	//PROD: 8001
	//DEV: 6969
	log.Fatal(app.Listen(":6969"))
	log.Println("Listening to port 6969")

}
