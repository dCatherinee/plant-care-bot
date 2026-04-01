package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	api                    *bot.Bot
	log                    *slog.Logger
	sendTextFn             func(ctx context.Context, chatID int64, text string) error
	sendTextWithKeyboardFn func(ctx context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error
}

func New(token string, logger *slog.Logger) (*Bot, error) {
	if token == "" {
		return nil, errors.New("telegram token is empty")
	}

	b, err := bot.New(
		token,
		bot.WithNotAsyncHandlers(),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("telegram polling error", "err", err)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	tgBot := &Bot{
		api: b,
		log: logger,
	}
	tgBot.sendTextFn = tgBot.sendTextMessage
	tgBot.sendTextWithKeyboardFn = tgBot.sendTextWithKeyboardMessage

	return tgBot, nil
}

func (b *Bot) Run(ctx context.Context) error {
	b.registerHandlers()

	b.log.Info("telegram bot started")
	b.api.Start(ctx)
	b.log.Info("telegram bot stopped")

	return ctx.Err()
}

func (b *Bot) sendText(ctx context.Context, chatID int64, text string) error {
	if b.sendTextFn != nil {
		return b.sendTextFn(ctx, chatID, text)
	}

	return b.sendTextMessage(ctx, chatID, text)
}

func (b *Bot) sendTextMessage(ctx context.Context, chatID int64, text string) error {
	_, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("send text message: %w", err)
	}

	return nil
}

func (b *Bot) sendTextWithKeyboard(ctx context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
	if b.sendTextWithKeyboardFn != nil {
		return b.sendTextWithKeyboardFn(ctx, chatID, text, keyboard)
	}

	return b.sendTextWithKeyboardMessage(ctx, chatID, text, keyboard)
}

func (b *Bot) sendTextWithKeyboardMessage(ctx context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
	_, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("send text message with keyboard: %w", err)
	}

	return nil
}

func (b *Bot) registerHandlers() {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, b.handleStart)

	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonPlants, bot.MatchTypeExact, b.handlePlants)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonCare, bot.MatchTypeExact, b.handleCare)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonReminders, bot.MatchTypeExact, b.handleReminders)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonSettings, bot.MatchTypeExact, b.handleSettings)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonHelp, bot.MatchTypeExact, b.handleHelp)
}
