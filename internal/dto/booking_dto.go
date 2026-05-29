package dto

type CreateBookingRequest struct {
	RoomName string `json:"room_name"`
	Date     string `json:"date"`
}