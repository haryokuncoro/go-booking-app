package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	JWTSecret     string
	JWTExpireHour int
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

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		
		JWTSecret: viper.GetString("JWT_SECRET"),
		JWTExpireHour: viper.GetInt("JWT_EXPIRE_HOUR"),
		
	}
}
