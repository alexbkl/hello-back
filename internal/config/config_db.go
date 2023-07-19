package config

import (
	"errors"
	"fmt"

	"github.com/Hello-Storage/hello-back/internal/db"
)

func (c *Config) ConnectDB() error {
	dbDsn := c.DatabaseDsn()

	log.Info("config: dbdsn is ", dbDsn)

	if dbDsn == "" {
		return errors.New("config: database DSN not specified")
	}

	dbconn := db.DbConn{
		Driver: "postgres",
		Dsn:    dbDsn,
	}

	dbconn.Open()
	db.SetDbConn(dbconn)

	return nil
}

// DatabaseDsn returns the database data source name (DSN).
func (c *Config) DatabaseDsn() string {
	dbDsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword)

	return dbDsn
}
