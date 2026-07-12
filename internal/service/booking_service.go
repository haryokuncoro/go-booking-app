package service

import (
	"booking-app/internal/cache"
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/lock"
	"booking-app/internal/logger"
	"booking-app/internal/repository"
	"booking-app/internal/worker"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
	"errors"
)

var ErrRoomAlreadyBooked = errors.New("room already booked for the requested date")

type BookingService interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, req dto.CreateBookingRequest) error
	GetBooking(ctx context.Context, id uuid.UUID) (*entity.Booking, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	userRepo    repository.UserRepository
	redis       *redis.Client
}

func NewBookingService(bookingRepo repository.BookingRepository, userRepo repository.UserRepository, redis *redis.Client) BookingService {

	return &bookingService{
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
		redis:       redis,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID uuid.UUID, req dto.CreateBookingRequest) error {

	date, err := time.Parse("2006-01-02", req.Date)

	if err != nil {
		logger.Log.Error(
			"create booking failed",

			zap.Error(err),
		)
		return err
	}

	roomLock := lock.GetRoomLock(req.RoomID)

	roomLock.Lock()

	defer roomLock.Unlock()

	existing, _ := s.bookingRepo.FindByRoomAndDate(ctx, req.RoomID, date)

	if existing != nil {
		logger.Log.Error(
			"room already booked",
			zap.Uint("room_id", req.RoomID),
			zap.String("date", req.Date),
		)
		return ErrRoomAlreadyBooked
	}

	booking := entity.Booking{
		UserID:      userID,
		RoomID:      req.RoomID,
		BookingDate: date,
		Status:      "CONFIRMED",
	}

	if err = s.bookingRepo.Create(ctx, &booking); err != nil {
		logger.Log.Error(
			"create booking failed",
			zap.Error(err),
		)
		return err
	}

	s.redis.Del(ctx, cache.BookingKey(booking.ID))

	user, err := s.userRepo.FindByID(ctx, userID)

	if err == nil {

		worker.EmailQueue <- worker.EmailJob{
			UserEmail: user.Email,
			RoomId:    booking.RoomID,
		}
	}
	logger.Log.Info(
		"booking created",
		zap.String("user_id", userID.String()),
		zap.Uint("room_id", booking.RoomID,
		),
	)
	return nil
}

func (s *bookingService) GetBooking(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	key := cache.BookingKey(id)

	cached, err := s.redis.Get(ctx, key).Result()

	if err == nil {

		var booking entity.Booking
		if unmarshalErr := json.Unmarshal([]byte(cached), &booking); unmarshalErr != nil {
			logger.Log.Error(
				"booking cache unmarshal failed",
				zap.String("booking_id", id.String()),
				zap.Error(unmarshalErr),
			)
		} else {
			logger.Log.Info(
				"booking cache hit",
				zap.String("booking_id", id.String()),
			)
			return &booking, nil
		}
	}

	booking, err := s.bookingRepo.FindByID(ctx, id)
	if err != nil {
		logger.Log.Error(
			"Get booking failed",
			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"booking cache miss",
		zap.String("booking_id", id.String()),
	)

	data, err := json.Marshal(booking)
	if err == nil {
		s.redis.Set(ctx, key, string(data), 5*time.Minute)
	}

	return booking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {

	return s.bookingRepo.FindByUserID(ctx, userID)
}
