package repository

import (
	"booking-app/internal/entity"
	"gorm.io/gorm"
	"context"
)


type BookingRepository interface {

	Create(
		ctx context.Context,
		booking *entity.Booking,
	) error

	FindByID(
		ctx context.Context,
		id uint,
	) (*entity.Booking, error)

	FindByUserID(
		ctx context.Context,
		userID uint,
	) ([]entity.Booking, error)
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
	id uint,
) (*entity.Booking, error) {

	var booking entity.Booking

	err := r.db.
		WithContext(ctx).
		First(&booking, id).
		Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindByUserID(
	ctx context.Context,
	userID uint,
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