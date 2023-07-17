package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	FilebaseAccessKey          string `mapstructure:"FILEBASE_ACCESS_KEY"`
	FilebaseSecretAcessKey     string `mapstructure:"FILEBASE_SECRET_ACCESS_KEY"`
	FilebasePinningAccessToken string `mapstructure:"FILEBASE_PINNING_ACCESS_TOKEN"`
}

var config Config

func Init() *Config {
	var err error

	config, err := load(".")

	if err != nil {
		logrus.Fatal("Could not load config: ", err)
	}

	return &config
}

func load(path string) (conf Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&conf)
	return
}
