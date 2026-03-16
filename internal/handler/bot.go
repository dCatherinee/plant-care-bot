package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

const (
	defaultHelpText = "Доступные команды:\n/start\nдобавить <название>\nсписок"
	addUsageText    = "Напиши: добавить <название растения>"
	addErrorText    = "Не удалось добавить растение."
	listErrorText   = "Не удалось получить список растений."
	unknownText     = "Я пока понимаю команды /start, добавить <название> и список."
)

type BotHandler struct {
	plants PlantUsecase
}

func NewBotHandler(plants PlantUsecase) *BotHandler {
	return &BotHandler{plants: plants}
}

func (h *BotHandler) HandleText(ctx context.Context, userID int64, text string) string {
	trimmed := strings.TrimSpace(text)
	normalized := strings.ToLower(trimmed)

	switch {
	case normalized == "/start":
		return defaultHelpText
	case normalized == "список":
		return h.handleList(ctx, userID)
	case strings.HasPrefix(normalized, "добавить"):
		return h.handleAdd(ctx, userID, trimmed)
	default:
		return unknownText
	}
}

func (h *BotHandler) handleAdd(ctx context.Context, userID int64, text string) string {
	name := strings.TrimSpace(strings.TrimPrefix(text, "добавить"))
	if name == text {
		lowerText := strings.ToLower(text)
		if strings.HasPrefix(lowerText, "добавить") {
			name = strings.TrimSpace(text[len("добавить"):])
		}
	}

	if name == "" {
		return addUsageText
	}

	plant, err := h.plants.AddPlant(ctx, userID, name)
	if err != nil {
		return addErrorText
	}

	return fmt.Sprintf("Добавила растение: %s", plant.Name)
}

func (h *BotHandler) handleList(ctx context.Context, userID int64) string {
	plants, err := h.plants.ListPlants(ctx, userID)
	if err != nil {
		return listErrorText
	}

	return formatPlantList(plants)
}

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
