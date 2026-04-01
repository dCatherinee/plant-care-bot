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
	api *bot.Bot
	log *slog.Logger
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

	return &Bot{
		api: b,
		log: logger,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.handleStart)

	b.log.Info("telegram bot started")
	b.api.Start(ctx)
	b.log.Info("telegram bot stopped")

	return ctx.Err()
}

func (b *Bot) handleStart(ctx context.Context, bt *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := bt.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет! Я бот для ухода за растениями 🌿",
	})
	if err != nil {
		b.log.Error("send /start response", "err", err)
	}
}
