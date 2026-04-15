package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
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

func (b *Bot) replyWithError(ctx context.Context, chatID int64, userErr error, keyboard models.ReplyKeyboardMarkup, logMessage string) {
	err := b.sendTextWithKeyboard(
		ctx,
		chatID,
		userMessageFromError(userErr),
		keyboard,
	)
	if err != nil {
		b.log.Error(logMessage, "err", err)
	}
}

func (b *Bot) pendingDeleteStore() *PendingDeleteStore {
	if b.pendingDeletes == nil {
		b.pendingDeletes = NewPendingDeleteStore()
	}

	return b.pendingDeletes
}

func (b *Bot) ensureTelegramUser(ctx context.Context, chatID, telegramUserID int64, keyboard models.ReplyKeyboardMarkup) (domain.User, bool) {
	user, err := b.users.EnsureUser(ctx, telegramUserID)
	if err != nil {
		b.replyWithError(ctx, chatID, err, keyboard, "send ensure user error")
		return domain.User{}, false
	}

	return user, true
}

func (b *Bot) handleListPlants(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramUserID := update.Message.From.ID

	user, ok := b.ensureTelegramUser(ctx, chatID, telegramUserID, plantsMenuKeyboard())
	if !ok {
		return
	}

	plants, err := b.plants.ListPlants(ctx, user.ID)
	if err != nil {
		b.replyWithError(ctx, chatID, err, plantsMenuKeyboard(), "send list plants error")
		return
	}

	err = b.sendTextWithKeyboard(
		ctx,
		chatID,
		formatPlantList(plants),
		plantsMenuKeyboard(),
	)
	if err != nil {
		b.log.Error("send plants list", "err", err)
	}
}

func (b *Bot) handleDeletePlant(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramUserID := update.Message.From.ID

	user, ok := b.ensureTelegramUser(ctx, chatID, telegramUserID, plantsMenuKeyboard())
	if !ok {
		return
	}

	plants, err := b.plants.ListPlants(ctx, user.ID)
	if err != nil {
		b.replyWithError(ctx, chatID, err, plantsMenuKeyboard(), "send delete plants list error")
		return
	}

	if len(plants) == 0 {
		err := b.sendTextWithKeyboard(ctx, chatID, "Список растений пуст.", plantsMenuKeyboard())
		if err != nil {
			b.log.Error("send empty plants list for delete", "err", err)
		}
		return
	}

	b.pendingDeleteStore().Clear(telegramUserID)
	inlineButtons := make([]models.InlineKeyboardButton, 0, len(plants))
	for _, plant := range plants {
		inlineButtons = append(inlineButtons, models.InlineKeyboardButton{
			Text:         plant.Name,
			CallbackData: callbackDeleteSelectPrefix + int64ToString(plant.ID),
		})
	}

	err = b.sendTextWithInlineKeyboard(
		ctx,
		chatID,
		"Выбери растение для удаления:",
		deletePlantsInlineKeyboard(inlineButtons),
	)
	if err != nil {
		b.log.Error("send delete plant prompt", "err", err)
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

func validatePlantName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return domain.ValidationError{
			Field:   "name",
			Problem: "is empty",
		}
	}

	if len([]rune(name)) > 50 {
		return domain.ValidationError{
			Field:   "name",
			Problem: "too long",
		}
	}

	return nil
}

func (b *Bot) handlePlantNameInput(ctx context.Context, chatID, userID int64, text string) {
	err := validatePlantName(text)
	if err != nil {
		b.replyWithError(ctx, chatID, err, cancelKeyboard(), "send empty plant name warning")
		return
	}

	user, ok := b.ensureTelegramUser(ctx, chatID, userID, cancelKeyboard())
	if !ok {
		return
	}

	plant, err := b.plants.AddPlant(ctx, user.ID, text)
	if err != nil {
		b.replyWithError(ctx, chatID, err, cancelKeyboard(), "send add plant error")
		return
	}

	b.states.Clear(userID)

	err = b.sendTextWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf("Растение \"%s\" добавлено 🌿", plant.Name),
		plantsMenuKeyboard(),
	)
	if err != nil {
		b.log.Error("send plant confirmation", "err", err)
	}
}

func (b *Bot) handleDeleteSelectCallback(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	if err := b.answerCallbackQuery(ctx, update.CallbackQuery.ID); err != nil {
		b.log.Error("answer delete select callback", "err", err)
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	telegramUserID := update.CallbackQuery.From.ID
	plantID, ok := parseCallbackPlantID(update.CallbackQuery.Data, callbackDeleteSelectPrefix)
	if !ok {
		return
	}

	user, ok := b.ensureTelegramUser(ctx, chatID, telegramUserID, plantsMenuKeyboard())
	if !ok {
		return
	}

	plant, err := b.plants.GetPlant(ctx, user.ID, plantID)
	if err != nil {
		b.replyWithError(ctx, chatID, err, plantsMenuKeyboard(), "send get plant for delete error")
		return
	}

	b.pendingDeleteStore().Set(telegramUserID, pendingDelete{
		userID:    user.ID,
		plantID:   plant.ID,
		plantName: plant.Name,
	})

	err = b.editMessageTextWithInlineKeyboard(
		ctx,
		chatID,
		messageID,
		fmt.Sprintf("Удалить растение \"%s\"?", plant.Name),
		deleteConfirmInlineKeyboard(plant.ID),
	)
	if err != nil {
		b.log.Error("send delete plant confirm prompt", "err", err)
	}
}

func (b *Bot) handleDeleteConfirmCallback(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	if err := b.answerCallbackQuery(ctx, update.CallbackQuery.ID); err != nil {
		b.log.Error("answer delete confirm callback", "err", err)
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	telegramUserID := update.CallbackQuery.From.ID

	plantID, ok := parseCallbackPlantID(update.CallbackQuery.Data, callbackDeleteConfirmPrefix)
	if !ok {
		return
	}

	pending, ok := b.pendingDeleteStore().Get(telegramUserID)
	if !ok || pending.plantID != plantID {
		err := b.editMessageTextWithInlineKeyboard(
			ctx,
			chatID,
			messageID,
			"Не удалось продолжить удаление. Попробуй заново.",
			models.InlineKeyboardMarkup{},
		)
		if err != nil {
			b.log.Error("send missing pending delete warning", "err", err)
		}
		return
	}

	err := b.plants.DeletePlant(ctx, pending.userID, pending.plantID)
	if err != nil {
		b.pendingDeleteStore().Clear(telegramUserID)
		b.replyWithError(ctx, chatID, err, plantsMenuKeyboard(), "send delete plant error")
		return
	}

	b.pendingDeleteStore().Clear(telegramUserID)

	err = b.editMessageTextWithInlineKeyboard(
		ctx,
		chatID,
		messageID,
		fmt.Sprintf("Растение \"%s\" удалено.", pending.plantName),
		models.InlineKeyboardMarkup{},
	)
	if err != nil {
		b.log.Error("send delete plant confirmation", "err", err)
	}
}

func (b *Bot) handleDeleteCancelCallback(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil || update.CallbackQuery.Message.Message == nil {
		return
	}

	if err := b.answerCallbackQuery(ctx, update.CallbackQuery.ID); err != nil {
		b.log.Error("answer delete cancel callback", "err", err)
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	telegramUserID := update.CallbackQuery.From.ID

	b.pendingDeleteStore().Clear(telegramUserID)

	err := b.editMessageTextWithInlineKeyboard(
		ctx,
		chatID,
		messageID,
		"Удаление отменено.",
		models.InlineKeyboardMarkup{},
	)
	if err != nil {
		b.log.Error("send delete cancel confirmation", "err", err)
	}
}

func (b *Bot) handleCancel(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.pendingDeleteStore().Clear(userID)
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

	b.pendingDeleteStore().Clear(userID)
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
