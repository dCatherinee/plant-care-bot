package telegram

import (
	"fmt"
	"strings"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
	"github.com/go-telegram/bot/models"
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

func formatCarePlantPrompt() string {
	return "Выбери растение для ухода:"
}

func formatCareActionsPrompt(plantName string) string {
	return fmt.Sprintf("Что отметить для \"%s\"?", plantName)
}

func formatCareMarkedMessage(plantName string, careType domain.CareKind) string {
	switch careType {
	case domain.CareKindWater:
		return fmt.Sprintf("Полив для %s отмечен.", plantName)
	case domain.CareKindFertilize:
		return fmt.Sprintf("Подкормка для %s отмечена.", plantName)
	default:
		return fmt.Sprintf("Уход для %s отмечен.", plantName)
	}
}

func carePlantButtons(plants []domain.Plant) []models.InlineKeyboardButton {
	buttons := make([]models.InlineKeyboardButton, 0, len(plants))
	for _, plant := range plants {
		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         plant.Name,
			CallbackData: callbackCareSelectPrefix + int64ToString(plant.ID),
		})
	}

	return buttons
}

func formatCareJournal(careType domain.CareKind, careEvents []domain.CareEvent, plantNames map[int64]string) string {
	title := "Журнал ухода:"
	switch careType {
	case domain.CareKindWater:
		title = "Журнал полива:"
	case domain.CareKindFertilize:
		title = "Журнал удобрений:"
	}

	if len(careEvents) == 0 {
		return title + "\nЗаписей пока нет."
	}

	lines := make([]string, 0, len(careEvents)+1)
	lines = append(lines, title)

	for i, careEvent := range careEvents {
		plantName := plantNames[careEvent.PlantID]
		if plantName == "" {
			plantName = fmt.Sprintf("Растение #%d", careEvent.PlantID)
		}

		lines = append(lines, fmt.Sprintf(
			"%d. %s — %s UTC",
			i+1,
			plantName,
			careEvent.OccurredAt.UTC().Format("02.01.2006 15:04"),
		))
	}

	return strings.Join(lines, "\n")
}
