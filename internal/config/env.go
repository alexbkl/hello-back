package config

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
)

type EnvVar struct {
	// App env
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`
	// token env
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	// Postgres
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
	// Filebase credential
	FilebaseBucket     string `mapstructure:"FILEBASE_BUCKET"`
	FilebaseAccessKey  string `mapstructure:"FILEBASE_ACCESS_KEY"`
	FilebaseSecretKey  string `mapstructure:"FILEBASE_SECRET_KEY"`
	FilebasePinningKey string `mapstructure:"FILEBASE_PINNING_KEY"`
}

var env EnvVar

func LoadEnv() (err error) {
	// skip load env when docker
	if os.Getenv("APP_PORT") == "" {
		err = godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	atd, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		return err
	}

	rtd, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		return err
	}

	env = EnvVar{
		// App env
		AppPort: os.Getenv("APP_PORT"),
		AppEnv:  os.Getenv("APP_ENV"),
		// token env
		TokenSymmetricKey:    os.Getenv("TOKEN_SYMMETRIC_KEY"),
		AccessTokenDuration:  atd,
		RefreshTokenDuration: rtd,
		// Postgres
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBPort:     os.Getenv("POSTGRES_PORT"),
		// Filebase credential
		FilebaseBucket:     os.Getenv("FILEBASE_BUCKET"),
		FilebaseAccessKey:  os.Getenv("FILEBASE_ACCESS_KEY"),
		FilebaseSecretKey:  os.Getenv("FILEBASE_SECRET_KEY"),
		FilebasePinningKey: os.Getenv("FILEBASE_PINNING_KEY"),
	}

	values := reflect.ValueOf(env)
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

func Env() EnvVar {
	return env
}
