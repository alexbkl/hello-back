package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type EnvVar struct {
	// App env
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`
	// Postgres
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
	// Filebase credential
	FilebaseAccessKey  string `mapstructure:"FILEBASE_ACCESS_KEY"`
	FilebaseSecretKey  string `mapstructure:"FILEBASE_SECRET_KEY"`
	FilebasePinningKey string `mapstructure:"FILEBASE_PINNING_KEY"`
}

var Env EnvVar

func LoadEnv() (err error) {
	err = godotenv.Load(".env")

	Env = EnvVar{
		AppPort:            os.Getenv("APP_PORT"),
		AppEnv:             os.Getenv("APP_ENV"),
		DBHost:             os.Getenv("POSTGRES_HOST"),
		DBName:             os.Getenv("POSTGRES_DB"),
		DBUser:             os.Getenv("POSTGRES_USER"),
		DBPassword:         os.Getenv("POSTGRES_PASSWORD"),
		DBPort:             os.Getenv("POSTGRES_PORT"),
		FilebaseAccessKey:  os.Getenv("FILEBASE_ACCESS_KEY"),
		FilebaseSecretKey:  os.Getenv("FILEBASE_SECRET_KEY"),
		FilebasePinningKey: os.Getenv("FILEBASE_PINNING_KEY"),
	}

	values := reflect.ValueOf(Env)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() == "" {
			return fmt.Errorf("config: %s is missing", types.Field(i).Name)
		}
	}

	if err != nil {
		return
	}

	return
}
