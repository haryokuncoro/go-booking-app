package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"size:100"`
	Email     string    `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
