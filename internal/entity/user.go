package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100"`
	Email     string    `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}