package telegram

import (
	"context"

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
	b.replyStub(ctx, update, `Раздел "Растения" пока в разработке 🌱`)
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
