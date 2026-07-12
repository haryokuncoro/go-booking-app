package service

import (
	"context"
	"errors"
	"testing"

	"booking-app/internal/entity"
)

func TestFindByID_Success(t *testing.T) {
	want := &entity.User{ID: 5, Name: "Frank", Email: "frank@example.com"}
	repo := &stubUserRepo{
		findByIDFn: func(ctx context.Context, id uint) (*entity.User, error) {
			if id != 5 {
				t.Errorf("expected id 5, got %d", id)
			}
			return want, nil
		},
	}
	svc := NewUserService(repo, testConfig())

	got, err := svc.FindByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != want {
		t.Errorf("expected user %+v, got %+v", want, got)
	}
}

func TestFindByID_RepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &stubUserRepo{
		findByIDFn: func(ctx context.Context, id uint) (*entity.User, error) {
			return nil, repoErr
		},
	}
	svc := NewUserService(repo, testConfig())

	got, err := svc.FindByID(context.Background(), 5)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil user on error, got %+v", got)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := &stubUserRepo{
		findByIDFn: func(ctx context.Context, id uint) (*entity.User, error) {
			return nil, nil
		},
	}
	svc := NewUserService(repo, testConfig())

	got, err := svc.FindByID(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
	if err.Error() != "User not found" {
		t.Errorf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil user, got %+v", got)
	}
}
