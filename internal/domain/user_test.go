package domain

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		user, err := NewUser(12345)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.TelegramUserID != 12345 {
			t.Fatalf("expected telegram user ID %d, got %d", 12345, user.TelegramUserID)
		}
	})

	t.Run("invalid_telegram_user_id", func(t *testing.T) {
		_, err := NewUser(0)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var validationErr ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}

		if validationErr.Field != "telegramUserID" {
			t.Fatalf("expected field %q, got %q", "telegramUserID", validationErr.Field)
		}

		if validationErr.Problem != "must be positive" {
			t.Fatalf("expected problem %q, got %q", "must be positive", validationErr.Problem)
		}

		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("expected ErrInvalidArgument, got %v", err)
		}
	})
}
