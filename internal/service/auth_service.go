package service

import (
	"errors"

	"booking-app/config"
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/repository"
	"booking-app/internal/utils"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(
		ctx context.Context,
		req dto.RegisterRequest,
	) error
	Login(
		ctx context.Context,
		req dto.LoginRequest,
	) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	cfg *config.Config,
) AuthService {

	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) error {

	existingUser, _ :=
		s.userRepo.FindByEmail(
			ctx,
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
		Name:  req.Name,
		Email: req.Email,
		Password: string(
			hashedPassword,
		),
	}

	return s.userRepo.Create(
		ctx,
		&user,
	)
}

func (s *authService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (string, error) {

	user, err :=
		s.userRepo.FindByEmail(
			ctx,
			req.Email,
		)

	if err != nil {
		return "", errors.New(
			"invalid email or password",
		)
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return "", errors.New(
			"invalid email or password",
		)
	}

	token, err :=
		utils.GenerateToken(
			user.ID,
			s.cfg.JWTSecret,
			s.cfg.JWTExpireHour,
		)

	if err != nil {
		return "", err
	}

	return token, nil
}
