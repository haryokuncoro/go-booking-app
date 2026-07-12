package entity

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid"`

	RoomID uint `gorm:"uniqueIndex:idx_room_date"`

	BookingDate time.Time `gorm:"uniqueIndex:idx_room_date"`

	Status string
}
