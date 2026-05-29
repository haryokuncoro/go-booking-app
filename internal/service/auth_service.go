package service

import (
	"errors"

	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(
		req dto.RegisterRequest,
	) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
) AuthService {

	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(
	req dto.RegisterRequest,
) error {

	existingUser, _ :=
		s.userRepo.FindByEmail(
			req.Email,
		)

	if existingUser != nil {
		return errors.New(
			"email already exists",
		)
	}

	hashedPassword, err :=
		bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

	if err != nil {
		return err
	}

	user := entity.User{
		Name: req.Name,
		Email: req.Email,
		Password: string(
			hashedPassword,
		),
	}

	return s.userRepo.Create(
		&user,
	)
}