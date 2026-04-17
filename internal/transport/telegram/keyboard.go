package telegram

import "github.com/go-telegram/bot/models"

const (
	buttonPlants    = "Растения"
	buttonCare      = "Уход"
	buttonReminders = "Напоминания"
	buttonSettings  = "Настройки"
	buttonHelp      = "Помощь"

	buttonAddPlant    = "Добавить растение"
	buttonListPlants  = "Список растений"
	buttonDeletePlant = "Удалить растение"
	buttonBackToMenu  = "Меню"
	buttonCancel      = "Отмена"

	callbackDeleteSelectPrefix  = "delete:select:"
	callbackDeleteConfirmPrefix = "delete:confirm:"
	callbackDeleteCancel        = "delete:cancel"
)

func mainMenuKeyboard() models.ReplyKeyboardMarkup {
	return models.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: buttonPlants},
				{Text: buttonCare},
			},
			{
				{Text: buttonReminders},
				{Text: buttonSettings},
			},
			{
				{Text: buttonHelp},
			},
		},
	}
}

func plantsMenuKeyboard() models.ReplyKeyboardMarkup {
	return models.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: buttonAddPlant},
				{Text: buttonListPlants},
			},
			{
				{Text: buttonDeletePlant},
			},
			{
				{Text: buttonBackToMenu},
			},
		},
	}
}

func cancelKeyboard() models.ReplyKeyboardMarkup {
	return models.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: buttonCancel},
				{Text: buttonBackToMenu},
			},
		},
	}
}

func deletePlantsInlineKeyboard(plants []models.InlineKeyboardButton) models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(plants)+1)
	for _, button := range plants {
		rows = append(rows, []models.InlineKeyboardButton{button})
	}

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         buttonCancel,
			CallbackData: callbackDeleteCancel,
		},
	})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func deleteConfirmInlineKeyboard(plantID int64) models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Удалить",
					CallbackData: callbackDeleteConfirmPrefix + int64ToString(plantID),
				},
				{
					Text:         buttonCancel,
					CallbackData: callbackDeleteCancel,
				},
			},
		},
	}
}

func emptyInlineKeyboard() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}
}
