package service

import (
	"booking-app/internal/repository"
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"time"
)

type BookingService interface {

	CreateBooking(
		userID uint,
		req dto.CreateBookingRequest,
	) error

	GetBooking(
		id uint,
	) (*entity.Booking, error)

	GetUserBookings(
		userID uint,
	) ([]entity.Booking, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
) BookingService {

	return &bookingService{
		bookingRepo: bookingRepo,
	}
}

func (s *bookingService) CreateBooking(
	userID uint,
	req dto.CreateBookingRequest,
) error {

	date, err :=
		time.Parse(
			"2006-01-02",
			req.Date,
		)

	if err != nil {
		return err
	}

	booking := entity.Booking{
		UserID:      userID,
		RoomName:    req.RoomName,
		BookingDate: date,
		Status:      "CONFIRMED",
	}

	return s.bookingRepo.Create(
		&booking,
	)
}

func (s *bookingService) GetBooking(
	id uint,
) (*entity.Booking, error) {

	return s.bookingRepo.FindByID(
		id,
	)
}

func (s *bookingService) GetUserBookings(
	userID uint,
) ([]entity.Booking, error) {

	return s.bookingRepo.
		FindByUserID(
			userID,
		)
}

