package dto

type CreateBookingRequest struct {
	RoomName string `json:"room_name" validate:"required,min=3,max=100"`
	Date string `json:"date" validate:"required"`
}