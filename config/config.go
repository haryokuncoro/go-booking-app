package config

import (
	"log"
	"os"

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

	RedisHost string
	RedisPort string
	
	JWTSecret     string
	JWTExpireHour int
}

func LoadConfig() *Config {

	viper.SetConfigFile(".env")

	// OS environment variables take precedence over values in the .env file.
	// This lets docker-compose `environment:` overrides (e.g. DB_HOST=postgres)
	// win over the localhost defaults baked into .env.
	viper.AutomaticEnv()

	// A missing .env file is not fatal: in containerized environments the
	// configuration comes entirely from environment variables.
	if err := viper.ReadInConfig(); err != nil && !os.IsNotExist(err) {
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

		RedisHost: viper.GetString("REDIS_HOST"),
		RedisPort: viper.GetString("REDIS_PORT"),
		
		JWTSecret: viper.GetString("JWT_SECRET"),
		JWTExpireHour: viper.GetInt("JWT_EXPIRE_HOUR"),
		
	}
}
