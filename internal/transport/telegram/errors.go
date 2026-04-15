package telegram

import (
	"errors"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func userMessageFromError(err error) string {
	switch {
	case errors.Is(err, domain.ErrPlantNameEmpty):
		return "Имя растения не должно быть пустым."

	case errors.Is(err, domain.ErrInvalidPlantName):
		return "Имя растения выглядит некорректно. Попробуй короче и без лишних символов."

	case errors.Is(err, domain.ErrPlantAlreadyExists):
		return "Растение с таким именем уже есть."

	default:
		return "Что-то пошло не так. Попробуй ещё раз позже."
	}
}
