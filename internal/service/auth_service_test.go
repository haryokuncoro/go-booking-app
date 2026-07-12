package service

import (
	"context"
	"errors"
	"testing"

	"booking-app/config"
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// stubUserRepo is a configurable UserRepository mock, driven by function fields
// so each test can tailor the behaviour it needs.
type stubUserRepo struct {
	createFn      func(ctx context.Context, user *entity.User) error
	findByEmailFn func(ctx context.Context, email string) (*entity.User, error)
	findByIDFn    func(ctx context.Context, id uint) (*entity.User, error)
}

func (s *stubUserRepo) Create(ctx context.Context, user *entity.User) error {
	if s.createFn != nil {
		return s.createFn(ctx, user)
	}
	return nil
}

func (s *stubUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if s.findByEmailFn != nil {
		return s.findByEmailFn(ctx, email)
	}
	return nil, errors.New("not found")
}

func (s *stubUserRepo) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret:     "test-secret",
		JWTExpireHour: 1,
	}
}

func TestRegister_Success(t *testing.T) {
	var created *entity.User
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
		createFn: func(ctx context.Context, user *entity.User) error {
			created = user
			return nil
		},
	}
	svc := NewAuthService(repo, testConfig())

	err := svc.Register(context.Background(), dto.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created == nil {
		t.Fatal("expected user to be created")
	}
	if created.Name != "Alice" || created.Email != "alice@example.com" {
		t.Errorf("unexpected user persisted: %+v", created)
	}
	if created.Password == "secret123" {
		t.Error("expected password to be hashed, got plaintext")
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(created.Password),
		[]byte("secret123"),
	); err != nil {
		t.Errorf("stored hash does not match original password: %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: 1, Email: email}, nil
		},
		createFn: func(ctx context.Context, user *entity.User) error {
			t.Fatal("Create should not be called when email exists")
			return nil
		},
	}
	svc := NewAuthService(repo, testConfig())

	err := svc.Register(context.Background(), dto.RegisterRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret123",
	})
	if err == nil {
		t.Fatal("expected error for existing email, got nil")
	}
	if err.Error() != "email already exists" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegister_CreateError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
		createFn: func(ctx context.Context, user *entity.User) error {
			return repoErr
		},
	}
	svc := NewAuthService(repo, testConfig())

	err := svc.Register(context.Background(), dto.RegisterRequest{
		Name:     "Carol",
		Email:    "carol@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	cfg := testConfig()
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:       42,
				Email:    email,
				Password: string(hashed),
			}, nil
		},
	}
	svc := NewAuthService(repo, cfg)

	token, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "dave@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a token, got empty string")
	}

	claims, err := utils.ParseToken(token, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("token should be parseable: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected UserID 42 in token, got %d", claims.UserID)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewAuthService(repo, testConfig())

	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "ghost@example.com",
		Password: "secret123",
	})
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
	if err.Error() != "invalid email or password" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:       7,
				Email:    email,
				Password: string(hashed),
			}, nil
		},
	}
	svc := NewAuthService(repo, testConfig())

	_, err := svc.Login(context.Background(), dto.LoginRequest{
		Email:    "eve@example.com",
		Password: "wrong-password",
	})
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
	if err.Error() != "invalid email or password" {
		t.Errorf("unexpected error: %v", err)
	}
}
