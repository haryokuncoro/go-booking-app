package service

import (
	"booking-app/internal/dto"
	"booking-app/internal/entity"
	"booking-app/internal/logger"
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	logger.Log = zap.NewNop()
	os.Exit(m.Run())
}

// --- mocks ---

type mockBookingRepo struct {
	mu       sync.Mutex
	bookings []*entity.Booking
}

func (m *mockBookingRepo) Create(ctx context.Context, b *entity.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b.ID = uuid.New()
	m.bookings = append(m.bookings, b)
	return nil
}

func (m *mockBookingRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, b := range m.bookings {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockBookingRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []entity.Booking
	for _, b := range m.bookings {
		if b.UserID == userID {
			result = append(result, *b)
		}
	}
	return result, nil
}

func (m *mockBookingRepo) FindByRoomAndDate(ctx context.Context, roomID uint, date time.Time) (*entity.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, b := range m.bookings {
		if b.RoomID == roomID && b.BookingDate.Equal(date) {
			return b, nil
		}
	}
	return nil, nil
}

type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, errors.New("not found")
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	// return error so worker.EmailQueue send is skipped
	return nil, errors.New("not found")
}

// --- helpers ---

func newTestService() BookingService {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	return NewBookingService(
		&mockBookingRepo{},
		&mockUserRepo{},
		redisClient,
	)
}

// --- tests ---

func TestCreateBooking_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.CreateBooking(ctx, uuid.New(), dto.CreateBookingRequest{
		RoomID: 1,
		Date:   "2026-06-01",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateBooking_InvalidDate(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.CreateBooking(ctx, uuid.New(), dto.CreateBookingRequest{
		RoomID: 1,
		Date:   "not-a-date",
	})

	if err == nil {
		t.Error("expected error for invalid date, got nil")
	}
}

func TestCreateBooking_RoomAlreadyBooked(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	req := dto.CreateBookingRequest{RoomID: 1, Date: "2026-06-01"}

	if err := svc.CreateBooking(ctx, uuid.New(), req); err != nil {
		t.Fatalf("first booking should succeed, got %v", err)
	}

	err := svc.CreateBooking(ctx, uuid.New(), req)
	if !errors.Is(err, ErrRoomAlreadyBooked) {
		t.Errorf("expected ErrRoomAlreadyBooked, got %v", err)
	}
}

func TestCreateBooking_ConcurrentSameRoomAndDate(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	req := dto.CreateBookingRequest{RoomID: 1, Date: "2026-06-01"}
	userID := uuid.New()

	var successCount atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := svc.CreateBooking(ctx, userID, req)
			if err == nil {
				successCount.Add(1)
			}
		}()
	}
	wg.Wait()

	if successCount.Load() != 1 {
		t.Errorf("expected 1 successful booking, got %d", successCount.Load())
	}
}

func TestGetBooking_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.GetBooking(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for missing booking, got nil")
	}
}

func TestGetUserBookings_Empty(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	bookings, err := svc.GetUserBookings(ctx, uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(bookings) != 0 {
		t.Errorf("expected empty slice, got %d bookings", len(bookings))
	}
}

func TestGetUserBookings_ReturnsOwned(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	user1 := uuid.New()
	user2 := uuid.New()

	_ = svc.CreateBooking(ctx, user1, dto.CreateBookingRequest{RoomID: 1, Date: "2026-06-01"})
	_ = svc.CreateBooking(ctx, user1, dto.CreateBookingRequest{RoomID: 2, Date: "2026-06-02"})
	_ = svc.CreateBooking(ctx, user2, dto.CreateBookingRequest{RoomID: 3, Date: "2026-06-03"})

	bookings, err := svc.GetUserBookings(ctx, user1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(bookings) != 2 {
		t.Errorf("expected 2 bookings for user 1, got %d", len(bookings))
	}
}
