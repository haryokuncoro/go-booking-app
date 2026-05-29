package service

import (
	"errors"
	"booking-app/config"
	"booking-app/internal/entity"
	"booking-app/internal/repository"
	"context"
)

type UserService interface{
	FindByID(
	ctx context.Context,
	id uint,
) (*entity.User, error)
}

type userService struct{
	userRepo 	repository.UserRepository
	cfg 		*config.Config
}

func NewUserService(
	userRepo repository.UserRepository, cfg *config.Config,
) UserService{
	return &userService{
		userRepo: userRepo,
		cfg : cfg,
	}
}

func(s *userService) FindByID(
	ctx context.Context,
	id uint,
) (*entity.User, error){
	existingUser, err :=
		s.userRepo.FindByID(
			ctx,
			id,
		)
	if err != nil{
		return nil, err
	}

	if existingUser == nil{
		return nil, errors.New("User not found")
	}

	return existingUser, nil
}
