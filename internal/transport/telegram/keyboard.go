package telegram

import "github.com/go-telegram/bot/models"

const (
	buttonPlants    = "Растения"
	buttonCare      = "Уход"
	buttonReminders = "Напоминания"
	buttonSettings  = "Настройки"
	buttonHelp      = "Помощь"
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
