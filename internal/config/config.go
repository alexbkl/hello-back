package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	FilebaseAccessKey          string `mapstructure:"FILEBASE_ACCESS_KEY"`
	FilebaseSecretAcessKey     string `mapstructure:"FILEBASE_SECRET_ACCESS_KEY"`
	FilebasePinningAccessToken string `mapstructure:"FILEBASE_PINNING_ACCESS_TOKEN"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath("./pkg/common/env")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
