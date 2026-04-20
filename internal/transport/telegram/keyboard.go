package telegram

import "github.com/go-telegram/bot/models"

const (
	buttonPlants    = "Растения"
	buttonCare      = "Уход"
	buttonReminders = "Напоминания"
	buttonSettings  = "Настройки"
	buttonHelp      = "Помощь"

	buttonAddPlant     = "Добавить растение"
	buttonListPlants   = "Список растений"
	buttonDeletePlant  = "Удалить растение"
	buttonBackToMenu   = "Меню"
	buttonCancel       = "Отмена"
	buttonCareMark     = "Отметить уход"
	buttonWaterLog     = "Журнал полива"
	buttonFertilizeLog = "Журнал удобрений"

	buttonMarkWater     = "Полил"
	buttonMarkFertilize = "Удобрил"
	buttonBackStep      = "Назад"

	callbackDeleteSelectPrefix  = "delete:select:"
	callbackDeleteConfirmPrefix = "delete:confirm:"
	callbackDeleteCancel        = "delete:cancel"
	callbackCareSelectPrefix    = "care:select:"
	callbackCareWaterPrefix     = "care:water:"
	callbackCareFertilizePrefix = "care:fertilize:"
	callbackCareBack            = "care:back"
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

func careMenuKeyboard() models.ReplyKeyboardMarkup {
	return models.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: buttonCareMark},
			},
			{
				{Text: buttonWaterLog},
				{Text: buttonFertilizeLog},
			},
			{
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

func carePlantsInlineKeyboard(plants []models.InlineKeyboardButton) models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(plants))
	for _, button := range plants {
		rows = append(rows, []models.InlineKeyboardButton{button})
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func careActionsInlineKeyboard(plantID int64) models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         buttonMarkWater,
					CallbackData: callbackCareWaterPrefix + int64ToString(plantID),
				},
				{
					Text:         buttonMarkFertilize,
					CallbackData: callbackCareFertilizePrefix + int64ToString(plantID),
				},
			},
			{
				{
					Text:         buttonBackStep,
					CallbackData: callbackCareBack,
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
