package telegram

import (
	"errors"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func userMessageFromError(err error) string {
	var validationErr domain.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return userMessageFromValidationError(validationErr)

	case errors.Is(err, domain.ErrPlantAlreadyExists):
		return "Растение с таким именем уже есть."

	default:
		return "Что-то пошло не так. Попробуй ещё раз позже."
	}
}

func userMessageFromValidationError(err domain.ValidationError) string {
	switch err.Field {
	case "name":
		if err.Problem == "is empty" {
			return "Имя растения не должно быть пустым."
		}

		return "Имя растения выглядит некорректно. Попробуй короче и без лишних символов."
	case "telegramUserID":
		return "Не удалось определить пользователя. Попробуй ещё раз позже."
	default:
		return "Проверь введённые данные и попробуй ещё раз."
	}
}
