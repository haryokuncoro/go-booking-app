package entity

import "time"

type Booking struct {
	ID uint `gorm:"primaryKey"`

	UserID uint

	RoomName string

	BookingDate time.Time

	Status string

	CreatedAt time.Time
	UpdatedAt time.Time
}