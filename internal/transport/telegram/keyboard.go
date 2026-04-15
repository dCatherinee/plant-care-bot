package telegram

import "github.com/go-telegram/bot/models"

const (
	buttonPlants    = "Растения"
	buttonCare      = "Уход"
	buttonReminders = "Напоминания"
	buttonSettings  = "Настройки"
	buttonHelp      = "Помощь"

	buttonAddPlant   = "Добавить растение"
	buttonBackToMenu = "Меню"
	buttonCancel     = "Отмена"
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
