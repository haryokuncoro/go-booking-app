package repository

import (
	"booking-app/internal/entity"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID)(*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
) UserRepository {

	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	user *entity.User,
) error {

	return r.db.
	WithContext(ctx).
	Create(user).Error
}

func (r *userRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*entity.User, error) {

	var user entity.User

	err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.User, error) {

	var user entity.User

	err := r.db.
		WithContext(ctx).
		First(&user, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
