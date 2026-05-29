package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppPort string
	AppEnv  string
}

func LoadConfig() *Config {

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		AppName: viper.GetString("APP_NAME"),
		AppPort: viper.GetString("APP_PORT"),
		AppEnv:  viper.GetString("APP_ENV"),
	}
}