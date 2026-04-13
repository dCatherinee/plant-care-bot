package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleStart(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := "Привет! Я бот для ухода за растениями 🌿\n\nВыбери раздел в меню ниже."

	err := b.sendTextWithKeyboard(ctx, update.Message.Chat.ID, text, mainMenuKeyboard())
	if err != nil {
		b.log.Error("send /start response", "err", err)
	}
}

func (b *Bot) handlePlants(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil {
		return
	}

	err := b.sendTextWithKeyboard(
		ctx,
		update.Message.Chat.ID,
		"Раздел растений 🌱\n\nВыбери действие.",
		plantsMenuKeyboard(),
	)
	if err != nil {
		b.log.Error("send plants menu", "err", err)
	}
}

func (b *Bot) handleCare(ctx context.Context, _ *bot.Bot, update *models.Update) {
	b.replyStub(ctx, update, `Раздел "Уход" пока в разработке 💧`)
}

func (b *Bot) handleReminders(ctx context.Context, _ *bot.Bot, update *models.Update) {
	b.replyStub(ctx, update, `Раздел "Напоминания" пока в разработке ⏰`)
}

func (b *Bot) handleSettings(ctx context.Context, _ *bot.Bot, update *models.Update) {
	b.replyStub(ctx, update, `Раздел "Настройки" пока в разработке ⚙️`)
}

func (b *Bot) handleHelp(ctx context.Context, _ *bot.Bot, update *models.Update) {
	b.replyStub(ctx, update, `Раздел "Помощь" пока в разработке ℹ️`)
}

func (b *Bot) replyStub(ctx context.Context, update *models.Update, text string) {
	if update == nil || update.Message == nil {
		return
	}

	err := b.sendText(ctx, update.Message.Chat.ID, text)
	if err != nil {
		b.log.Error("send stub response", "err", err)
	}
}

func (b *Bot) handleAddPlant(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	b.states.Set(userID, StateWaitingPlantName)

	err := b.sendTextWithKeyboard(
		ctx,
		update.Message.Chat.ID,
		"Введи имя растения.\n\nЧтобы выйти, нажми «Отмена» или «Меню».",
		cancelKeyboard(),
	)
	if err != nil {
		b.log.Error("send add plant prompt", "err", err)
	}
}

func (b *Bot) handleTextByState(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	switch b.states.Get(userID) {
	case StateWaitingPlantName:
		b.handlePlantNameInput(ctx, chatID, userID, text)
	default:
		b.sendTextMessage(ctx, chatID, "Неизвестная команда")
	}
}

func (b *Bot) handlePlantNameInput(ctx context.Context, chatID, userID int64, text string) {
	name := strings.TrimSpace(text)
	if name == "" {
		err := b.sendTextWithKeyboard(
			ctx,
			chatID,
			"Имя растения не должно быть пустым. Введи название ещё раз.",
			cancelKeyboard(),
		)
		if err != nil {
			b.log.Error("send empty plant name warning", "err", err)
		}
		return
	}

	b.states.Clear(userID)

	err := b.sendTextWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf("Растение \"%s\" добавлено 🌿", name),
		plantsMenuKeyboard(),
	)
	if err != nil {
		b.log.Error("send plant confirmation", "err", err)
	}
}

func (b *Bot) handleCancel(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.states.Clear(userID)

	err := b.sendTextWithKeyboard(
		ctx,
		chatID,
		"Действие отменено.",
		plantsMenuKeyboard(),
	)
	if err != nil {
		b.log.Error("send cancel response", "err", err)
	}
}

func (b *Bot) handleBackToMenu(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.states.Clear(userID)

	err := b.sendTextWithKeyboard(
		ctx,
		chatID,
		"Возвращаю в главное меню.",
		mainMenuKeyboard(),
	)
	if err != nil {
		b.log.Error("send back response", "err", err)
	}
}
