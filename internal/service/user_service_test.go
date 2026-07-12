package service

import (
	"context"
	"errors"
	"testing"

	"booking-app/internal/entity"

	"github.com/google/uuid"
)

func TestFindByID_Success(t *testing.T) {
	wantID := uuid.New()
	want := &entity.User{ID: wantID, Name: "Frank", Email: "frank@example.com"}
	repo := &stubUserRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			if id != wantID {
				t.Errorf("expected id %s, got %s", wantID, id)
			}
			return want, nil
		},
	}
	svc := NewUserService(repo, testConfig())

	got, err := svc.FindByID(context.Background(), wantID)
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
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			return nil, repoErr
		},
	}
	svc := NewUserService(repo, testConfig())

	got, err := svc.FindByID(context.Background(), uuid.New())
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil user on error, got %+v", got)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := &stubUserRepo{
		findByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			return nil, nil
		},
	}
	svc := NewUserService(repo, testConfig())

	got, err := svc.FindByID(context.Background(), uuid.New())
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
