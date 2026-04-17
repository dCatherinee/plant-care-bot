package telegram

import (
	"errors"
	"testing"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func TestUserMessageFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "empty plant name",
			err: domain.ValidationError{
				Field:   "name",
				Problem: "is empty",
			},
			want: "Имя растения не должно быть пустым.",
		},
		{
			name: "invalid plant name",
			err: domain.ValidationError{
				Field:   "name",
				Problem: "has invalid format",
			},
			want: "Имя растения выглядит некорректно. Попробуй короче и без лишних символов.",
		},
		{
			name: "invalid telegram user id",
			err: domain.ValidationError{
				Field:   "telegramUserID",
				Problem: "must be positive",
			},
			want: "Не удалось определить пользователя. Попробуй ещё раз позже.",
		},
		{
			name: "plant already exists",
			err:  domain.ErrPlantAlreadyExists,
			want: "Растение с таким именем уже есть.",
		},
		{
			name: "not found",
			err:  domain.ErrNotFound,
			want: "Растение не найдено.",
		},
		{
			name: "unexpected error",
			err:  errors.New("db failed"),
			want: "Что-то пошло не так. Попробуй ещё раз позже.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := userMessageFromError(tc.err); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
