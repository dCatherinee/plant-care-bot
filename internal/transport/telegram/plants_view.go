package telegram

import (
	"fmt"
	"strings"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func formatPlantList(plants []domain.Plant) string {
	if len(plants) == 0 {
		return "Список растений пуст."
	}

	lines := make([]string, 0, len(plants)+1)
	lines = append(lines, "Твои растения:")

	for i, plant := range plants {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, plant.Name))
	}

	return strings.Join(lines, "\n")
}
