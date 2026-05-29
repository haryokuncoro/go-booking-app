package service

import (
	"booking-app/internal/cache"
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/logger"
	"booking-app/internal/repository"
	"booking-app/internal/worker"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
	redis       *redis.Client
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	userRepo repository.UserRepository,
	redis *redis.Client,
) BookingService {

	return &bookingService{
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
		redis:       redis,
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
		logger.Log.Error(
			"create booking failed",

			zap.Error(err),
		)
		return err
	}

	booking := entity.Booking{
		UserID:      userID,
		RoomName:    req.RoomName,
		BookingDate: date,
		Status:      "CONFIRMED",
	}

	if err = s.bookingRepo.Create(&booking); err != nil {
		logger.Log.Error(
			"create booking failed",

			zap.Error(err),
		)
		return err
	}

	ctx := context.Background()
	s.redis.Del(ctx, cache.BookingKey(booking.ID))

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
	logger.Log.Info(
		"booking created",

		zap.Uint(
			"user_id",
			userID,
		),

		zap.String(
			"room_name",
			booking.RoomName,
		),
	)
	return nil
}

func (s *bookingService) GetBooking(
	id uint,
) (*entity.Booking, error) {
	ctx := context.Background()
	key := cache.BookingKey(id)

	cached, err :=
		s.redis.Get(
			ctx,
			key,
		).Result()

	if err == nil {

		var booking entity.Booking

		json.Unmarshal(
			[]byte(cached),
			&booking,
		)

		logger.Log.Info(
			"booking cache hit",

			zap.Uint(
				"booking_id",
				id,
			),
		)

		return &booking, nil
	}

	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		logger.Log.Error(
			"Get booking failed",

			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"booking cache miss",

		zap.Uint(
			"booking_id",
			id,
		),
	)

	data, err := json.Marshal(booking)
	if err == nil {
		s.redis.Set(ctx, key, string(data), 5*time.Minute)
	}

	return booking, nil
}

func (s *bookingService) GetUserBookings(
	userID uint,
) ([]entity.Booking, error) {

	return s.bookingRepo.
		FindByUserID(
			userID,
		)
}
