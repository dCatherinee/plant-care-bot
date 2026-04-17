package telegram

import (
	"testing"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func TestFormatPlantList(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := formatPlantList(nil)

		if got != "Список растений пуст." {
			t.Fatalf("expected empty list text, got %q", got)
		}
	})

	t.Run("with_plants", func(t *testing.T) {
		got := formatPlantList([]domain.Plant{
			{ID: 1, UserID: 10, Name: "Monstera"},
			{ID: 2, UserID: 10, Name: "Cactus"},
		})

		want := "Твои растения:\n1. Monstera\n2. Cactus"
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}
