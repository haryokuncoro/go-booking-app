package repository

import (
	"booking-app/internal/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BookingRepository interface {
	Create(
		ctx context.Context,
		booking *entity.Booking,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.Booking, error)

	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]entity.Booking, error)
	FindByRoomAndDate(
		ctx context.Context,
		roomID uint,
		date time.Time,
	) (*entity.Booking, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(
	db *gorm.DB,
) BookingRepository {

	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) Create(
	ctx context.Context,
	booking *entity.Booking,
) error {

	return r.db.
		WithContext(ctx).
		Create(booking).
		Error
}

func (r *bookingRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Booking, error) {

	var booking entity.Booking

	err := r.db.
		WithContext(ctx).
		First(&booking, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Booking, error) {

	var bookings []entity.Booking

	err := r.db.
		WithContext(ctx).
		Where(
			"user_id = ?",
			userID,
		).
		Find(&bookings).
		Error

	return bookings, err
}

func (r *bookingRepository) FindByRoomAndDate(
	ctx context.Context,
	roomID uint,
	date time.Time,
) (*entity.Booking, error) {

	var booking entity.Booking

	err := r.db.
		WithContext(ctx).
		Where(
			"room_id = ? AND booking_date = ?",
			roomID,
			date,
		).
		First(&booking).
		Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}