package dto

type CreateBookingRequest struct {
	RoomID uint `json:"room_id" validate:"required"`
	Date string `json:"date" validate:"required"`
}