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
	api                                 *bot.Bot
	log                                 *slog.Logger
	plants                              PlantUsecase
	users                               UserUsecase
	states                              *StateStore
	pendingDeletes                      *PendingDeleteStore
	sendTextFn                          func(ctx context.Context, chatID int64, text string) error
	sendTextWithKeyboardFn              func(ctx context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error
	sendTextWithInlineKeyboardFn        func(ctx context.Context, chatID int64, text string, keyboard models.InlineKeyboardMarkup) error
	editMessageTextWithInlineKeyboardFn func(ctx context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error
	answerCallbackQueryFn               func(ctx context.Context, callbackQueryID string) error
}

func New(token string, logger *slog.Logger, plants PlantUsecase, users UserUsecase) (*Bot, error) {
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
		api:            b,
		log:            logger,
		plants:         plants,
		users:          users,
		states:         NewStateStore(),
		pendingDeletes: NewPendingDeleteStore(),
	}
	tgBot.sendTextFn = tgBot.sendTextMessage
	tgBot.sendTextWithKeyboardFn = tgBot.sendTextWithKeyboardMessage
	tgBot.sendTextWithInlineKeyboardFn = tgBot.sendTextWithInlineKeyboardMessage
	tgBot.editMessageTextWithInlineKeyboardFn = tgBot.editMessageTextWithInlineKeyboardMessage
	tgBot.answerCallbackQueryFn = tgBot.answerCallbackQuery

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

func (b *Bot) sendTextWithInlineKeyboard(ctx context.Context, chatID int64, text string, keyboard models.InlineKeyboardMarkup) error {
	if b.sendTextWithInlineKeyboardFn != nil {
		return b.sendTextWithInlineKeyboardFn(ctx, chatID, text, keyboard)
	}

	return b.sendTextWithInlineKeyboardMessage(ctx, chatID, text, keyboard)
}

func (b *Bot) sendTextWithInlineKeyboardMessage(ctx context.Context, chatID int64, text string, keyboard models.InlineKeyboardMarkup) error {
	_, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("send text message with inline keyboard: %w", err)
	}

	return nil
}

func (b *Bot) editMessageTextWithInlineKeyboard(ctx context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
	if b.editMessageTextWithInlineKeyboardFn != nil {
		return b.editMessageTextWithInlineKeyboardFn(ctx, chatID, messageID, text, keyboard)
	}

	return b.editMessageTextWithInlineKeyboardMessage(ctx, chatID, messageID, text, keyboard)
}

func (b *Bot) editMessageTextWithInlineKeyboardMessage(ctx context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
	_, err := b.api.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("edit message text with inline keyboard: %w", err)
	}

	return nil
}

func (b *Bot) answerCallbackQuery(ctx context.Context, callbackQueryID string) error {
	if b.answerCallbackQueryFn != nil {
		return b.answerCallbackQueryFn(ctx, callbackQueryID)
	}

	_, err := b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQueryID,
	})
	if err != nil {
		return fmt.Errorf("answer callback query: %w", err)
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

	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonAddPlant, bot.MatchTypeExact, b.handleAddPlant)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonListPlants, bot.MatchTypeExact, b.handleListPlants)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonDeletePlant, bot.MatchTypeExact, b.handleDeletePlant)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonBackToMenu, bot.MatchTypeExact, b.handleBackToMenu)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, buttonCancel, bot.MatchTypeExact, b.handleCancel)
	b.api.RegisterHandler(bot.HandlerTypeCallbackQueryData, callbackDeleteCancel, bot.MatchTypeExact, b.handleDeleteCancelCallback)
	b.api.RegisterHandler(bot.HandlerTypeCallbackQueryData, callbackDeleteSelectPrefix, bot.MatchTypePrefix, b.handleDeleteSelectCallback)
	b.api.RegisterHandler(bot.HandlerTypeCallbackQueryData, callbackDeleteConfirmPrefix, bot.MatchTypePrefix, b.handleDeleteConfirmCallback)

	b.api.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.Message != nil
	}, b.handleTextByState)
}
