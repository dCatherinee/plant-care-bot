package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type fakeUserRepo struct {
	ensureUserFn func(ctx context.Context, telegramUserID int64) (domain.User, error)
}

func (f *fakeUserRepo) EnsureUser(ctx context.Context, telegramUserID int64) (domain.User, error) {
	if f.ensureUserFn != nil {
		return f.ensureUserFn(ctx, telegramUserID)
	}

	return domain.User{}, nil
}

func TestUserServiceEnsureUser(t *testing.T) {
	createdAt := time.Date(2026, time.April, 15, 10, 0, 0, 0, time.UTC)
	repo := &fakeUserRepo{
		ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
			if telegramUserID != 12345 {
				t.Fatalf("expected telegram user ID %d, got %d", 12345, telegramUserID)
			}

			return domain.User{
				ID:             7,
				TelegramUserID: telegramUserID,
				CreatedAt:      createdAt,
			}, nil
		},
	}

	service := NewUserService(repo)

	user, err := service.EnsureUser(context.Background(), 12345)
	mustNoErr(t, err)

	if user.ID != 7 {
		t.Fatalf("expected user ID %d, got %d", 7, user.ID)
	}

	if user.TelegramUserID != 12345 {
		t.Fatalf("expected telegram user ID %d, got %d", 12345, user.TelegramUserID)
	}

	if user.CreatedAt != createdAt {
		t.Fatalf("expected createdAt %v, got %v", createdAt, user.CreatedAt)
	}
}

func TestUserServiceEnsureUserValidationError(t *testing.T) {
	repoCalled := false
	repo := &fakeUserRepo{
		ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
			repoCalled = true
			return domain.User{}, nil
		},
	}

	service := NewUserService(repo)

	_, err := service.EnsureUser(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if repoCalled {
		t.Fatal("repo should not be called on invalid input")
	}

	if !errors.Is(err, domain.ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestUserServiceEnsureUserCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	repoCalled := false
	repo := &fakeUserRepo{
		ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
			repoCalled = true
			return domain.User{}, nil
		},
	}

	service := NewUserService(repo)

	_, err := service.EnsureUser(ctx, 12345)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	if repoCalled {
		t.Fatal("repo should not be called when context is canceled")
	}
}
