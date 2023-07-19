package config

import (
	"github.com/Hello-Storage/hello-back/entities"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB
var DATABASE_URI string = "host=localhost user=postgres password=12345 dbname=metamask port=5432 sslmode=disable TimeZone=Europe/Madrid"

func Connect() (*gorm.DB, error) {
	var err error

	filterLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
		},
	)

	Database, err := gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 filterLogger,
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
	Database.AutoMigrate(&entities.PublishedFile{})
	Database.AutoMigrate(&entities.FileSharedState{})

	// Drop the old foreign key constraint to be able to delete the published file
	Database.Exec("ALTER TABLE file_shared_states DROP CONSTRAINT IF EXISTS fk_file_shared_states_published_file;")

	// Add the new foreign key constraint with ON DELETE SET NULL
	Database.Exec("ALTER TABLE file_shared_states ADD CONSTRAINT fk_file_shared_states_published_file FOREIGN KEY (published_file_id) REFERENCES published_files(id) ON DELETE SET NULL;")

	return Database, nil
}
