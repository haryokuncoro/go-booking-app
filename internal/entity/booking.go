package entity

import "time"

type Booking struct {
	ID uint `gorm:"primaryKey"`

	UserID uint

	RoomID uint `gorm:"uniqueIndex:idx_room_date"`

	BookingDate time.Time `gorm:"uniqueIndex:idx_room_date"`

	Status string
}