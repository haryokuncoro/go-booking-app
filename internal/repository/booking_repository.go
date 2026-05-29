package repository

import (
	"booking-app/internal/entity"
	"gorm.io/gorm"
)


type BookingRepository interface {

	Create(
		booking *entity.Booking,
	) error

	FindByID(
		id uint,
	) (*entity.Booking, error)

	FindByUserID(
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
	booking *entity.Booking,
) error {

	return r.db.
		Create(booking).
		Error
}

func (r *bookingRepository) FindByID(
	id uint,
) (*entity.Booking, error) {

	var booking entity.Booking

	err := r.db.
		First(&booking, id).
		Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindByUserID(
	userID uint,
) ([]entity.Booking, error) {

	var bookings []entity.Booking

	err := r.db.
		Where(
			"user_id = ?",
			userID,
		).
		Find(&bookings).
		Error

	return bookings, err
}