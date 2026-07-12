package dto

import "github.com/google/uuid"

type BookingResponse struct {
	ID uuid.UUID `json:"id"`

	RoomName string `json:"room_name"`

	Status string `json:"status"`

	Date string `json:"date"`
}