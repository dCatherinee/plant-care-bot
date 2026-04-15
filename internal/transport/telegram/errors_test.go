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
			err:  domain.ErrPlantNameEmpty,
			want: "Имя растения не должно быть пустым.",
		},
		{
			name: "invalid plant name",
			err:  domain.ErrInvalidPlantName,
			want: "Имя растения выглядит некорректно. Попробуй короче и без лишних символов.",
		},
		{
			name: "plant already exists",
			err:  domain.ErrPlantAlreadyExists,
			want: "Растение с таким именем уже есть.",
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
