package dto

type BookingResponse struct {
	ID uint `json:"id"`

	RoomName string `json:"room_name"`

	Status string `json:"status"`

	Date string `json:"date"`
}