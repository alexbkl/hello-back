package config

import (
	"errors"
	"fmt"

	"github.com/Hello-Storage/hello-back/internal/db"
	"github.com/Hello-Storage/hello-back/internal/entity"
	"github.com/Hello-Storage/hello-back/internal/migrate"
)

func (c *Config) ConnectDB() error {
	dbDsn := c.DatabaseDsn()

	if dbDsn == "" {
		return errors.New("config: database DSN not specified")
	}

	dbconn := db.DbConn{
		Driver: "postgres",
		Dsn:    dbDsn,
	}

	dbconn.Open()
	db.SetDbProvider(dbconn)

	return nil
}

// MigrateDb initializes the database and migrates the schema if needed.
func (c *Config) MigrateDb(runFailed bool, ids []string) {

	entity.InitDb(migrate.Opt(true, runFailed, ids))

	go entity.Error{}.LogEvents()
}

// DatabaseDsn returns the database data source name (DSN).
func (c *Config) DatabaseDsn() string {
	dbDsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword)

	return dbDsn
}
