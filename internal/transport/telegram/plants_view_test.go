package telegram

import (
	"strings"
	"testing"
	"time"

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

func TestFormatCareJournal(t *testing.T) {
	got := formatCareJournal(
		domain.CareKindWater,
		[]domain.CareEvent{
			{
				ID:         1,
				PlantID:    10,
				Kind:       domain.CareKindWater,
				OccurredAt: time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC),
			},
		},
		map[int64]string{
			10: "Monstera",
		},
	)

	if got == "" {
		t.Fatal("expected non-empty care journal text")
	}

	if !strings.HasPrefix(got, "Журнал полива:") {
		t.Fatalf("expected water journal title, got %q", got)
	}
}
