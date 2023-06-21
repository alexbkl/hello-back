package config

import (
	"meta-go-api/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB
var DATABASE_URI string = "host=localhost user=postgres password=12345 dbname=metamask port=5432 sslmode=disable TimeZone=Europe/Madrid"


func Connect() (*gorm.DB, error) {
	var err error

	Database, err := gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
        PrepareStmt:            true,
	})

	if err != nil {
		panic(err)
	} else {
		println("Database connected successfully")
	}

	// Migrate the schema
	Database.AutoMigrate(&entities.Dog{})
	Database.AutoMigrate(&entities.User{})
	Database.AutoMigrate(&entities.File{})
	Database.AutoMigrate(&entities.Email{})

	return Database, nil
}
