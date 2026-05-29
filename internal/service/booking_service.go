package service

import (
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/repository"
	"booking-app/internal/worker"
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
	userRepo    repository.UserRepository

}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	userRepo repository.UserRepository,
) BookingService {

	return &bookingService{
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
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

	if err = s.bookingRepo.Create(&booking); err != nil {
		return err
	}

	user, err :=
	s.userRepo.FindByID(
		userID,
	)

	if err == nil {

		worker.EmailQueue <- worker.EmailJob{
			UserEmail: user.Email,
			RoomName:  booking.RoomName,
		}
	}
	return nil
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

