package database

import (
	"booking-app/config"
	"booking-app/internal/entity"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(
	cfg *config.Config,
) *gorm.DB {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal(err)
	}
	
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Booking{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
