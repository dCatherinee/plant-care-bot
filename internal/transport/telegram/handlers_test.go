package telegram

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
	"github.com/go-telegram/bot/models"
)

type plantUsecaseStub struct {
	addPlantFn        func(ctx context.Context, userID int64, name string) (domain.Plant, error)
	listPlantsFn      func(ctx context.Context, userID int64) ([]domain.Plant, error)
	getPlantFn        func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error)
	updatePlantNameFn func(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error)
	deletePlantFn     func(ctx context.Context, userID int64, plantID int64) error
}

func (s plantUsecaseStub) AddPlant(ctx context.Context, userID int64, name string) (domain.Plant, error) {
	if s.addPlantFn != nil {
		return s.addPlantFn(ctx, userID, name)
	}

	return domain.Plant{}, nil
}

func (s plantUsecaseStub) ListPlants(ctx context.Context, userID int64) ([]domain.Plant, error) {
	if s.listPlantsFn != nil {
		return s.listPlantsFn(ctx, userID)
	}

	return nil, nil
}

func (s plantUsecaseStub) GetPlant(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
	if s.getPlantFn != nil {
		return s.getPlantFn(ctx, userID, plantID)
	}

	return domain.Plant{}, nil
}

func (s plantUsecaseStub) UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
	if s.updatePlantNameFn != nil {
		return s.updatePlantNameFn(ctx, userID, plantID, name)
	}

	return domain.Plant{}, nil
}

func (s plantUsecaseStub) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	if s.deletePlantFn != nil {
		return s.deletePlantFn(ctx, userID, plantID)
	}

	return nil
}

type userUsecaseStub struct {
	ensureUserFn func(ctx context.Context, telegramUserID int64) (domain.User, error)
}

func (s userUsecaseStub) EnsureUser(ctx context.Context, telegramUserID int64) (domain.User, error) {
	if s.ensureUserFn != nil {
		return s.ensureUserFn(ctx, telegramUserID)
	}

	return domain.User{}, nil
}

type careUsecaseStub struct {
	markWaterFn                  func(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error)
	markFertilizeFn              func(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error)
	listRecentCareEventsByTypeFn func(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error)
}

func (s careUsecaseStub) MarkWater(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
	if s.markWaterFn != nil {
		return s.markWaterFn(ctx, userID, plantID)
	}

	return domain.CareEvent{}, nil
}

func (s careUsecaseStub) MarkFertilize(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
	if s.markFertilizeFn != nil {
		return s.markFertilizeFn(ctx, userID, plantID)
	}

	return domain.CareEvent{}, nil
}

func (s careUsecaseStub) ListRecentCareEventsByType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
	if s.listRecentCareEventsByTypeFn != nil {
		return s.listRecentCareEventsByTypeFn(ctx, userID, eventType, limit)
	}

	return nil, nil
}

func TestHandleStartSendsWelcomeTextWithKeyboard(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleStart(context.Background(), nil, testUpdate(42))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	wantText := "Привет! Я бот для ухода за растениями 🌿\n\nВыбери раздел в меню ниже."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, mainMenuKeyboard()) {
		t.Fatalf("expected main menu keyboard, got %#v", gotKeyboard)
	}
}

func TestHandleStartIgnoresUpdateWithoutMessage(t *testing.T) {
	called := false

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		sendTextWithKeyboardFn: func(_ context.Context, _ int64, _ string, _ models.ReplyKeyboardMarkup) error {
			called = true
			return nil
		},
	}

	b.handleStart(context.Background(), nil, &models.Update{})

	if called {
		t.Fatal("expected sender not to be called")
	}
}

func TestStubHandlersSendExpectedText(t *testing.T) {
	tests := []struct {
		name    string
		handler func(ctx context.Context, bt *Bot, update *models.Update)
		want    string
	}{
		{
			name: "reminders",
			handler: func(ctx context.Context, bt *Bot, update *models.Update) {
				bt.handleReminders(ctx, nil, update)
			},
			want: `Раздел "Напоминания" пока в разработке ⏰`,
		},
		{
			name: "settings",
			handler: func(ctx context.Context, bt *Bot, update *models.Update) {
				bt.handleSettings(ctx, nil, update)
			},
			want: `Раздел "Настройки" пока в разработке ⚙️`,
		},
		{
			name: "help",
			handler: func(ctx context.Context, bt *Bot, update *models.Update) {
				bt.handleHelp(ctx, nil, update)
			},
			want: `Раздел "Помощь" пока в разработке ℹ️`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotChatID int64
			var gotText string

			b := &Bot{
				log: slog.New(slog.NewTextHandler(io.Discard, nil)),
				sendTextFn: func(_ context.Context, chatID int64, text string) error {
					gotChatID = chatID
					gotText = text
					return nil
				},
			}

			tt.handler(context.Background(), b, testUpdate(99))

			if gotChatID != 99 {
				t.Fatalf("expected chat ID %d, got %d", 99, gotChatID)
			}

			if gotText != tt.want {
				t.Fatalf("expected text %q, got %q", tt.want, gotText)
			}
		})
	}
}

func TestHandleCareStartsFlow(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleCare(context.Background(), nil, testUpdateFromUser(42, 1001, buttonCare))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	wantText := "Раздел ухода 💧\n\nВыбери действие."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, careMenuKeyboard()) {
		t.Fatalf("expected care menu keyboard, got %#v", gotKeyboard)
	}
}

func TestHandleCareMarkStartsInlineFlow(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.InlineKeyboardMarkup

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{ID: 77, TelegramUserID: telegramUserID}, nil
			},
		},
		plants: plantUsecaseStub{
			listPlantsFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
				return []domain.Plant{
					{ID: 1, UserID: userID, Name: "Monstera"},
					{ID: 2, UserID: userID, Name: "Cactus"},
				}, nil
			},
		},
		sendTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.InlineKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleCareMark(context.Background(), nil, testUpdateFromUser(42, 1001, buttonCareMark))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	if gotText != formatCarePlantPrompt() {
		t.Fatalf("expected text %q, got %q", formatCarePlantPrompt(), gotText)
	}

	if gotKeyboard.InlineKeyboard[0][0].CallbackData != callbackCareSelectPrefix+"1" {
		t.Fatalf("unexpected first care button: %#v", gotKeyboard.InlineKeyboard[0][0])
	}
}

func TestHandleWaterLogShowsRecentEvents(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup

	eventTime := time.Date(2026, 4, 17, 10, 30, 0, 0, time.UTC)

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{ID: 77, TelegramUserID: telegramUserID}, nil
			},
		},
		plants: plantUsecaseStub{
			listPlantsFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
				return []domain.Plant{
					{ID: 1, UserID: userID, Name: "Monstera"},
				}, nil
			},
		},
		care: careUsecaseStub{
			listRecentCareEventsByTypeFn: func(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
				if userID != 77 {
					t.Fatalf("expected user ID 77, got %d", userID)
				}
				if eventType != domain.CareKindWater {
					t.Fatalf("expected water kind, got %q", eventType)
				}
				if limit != 10 {
					t.Fatalf("expected limit 10, got %d", limit)
				}

				return []domain.CareEvent{
					{ID: 1, PlantID: 1, Kind: domain.CareKindWater, OccurredAt: eventTime},
				}, nil
			},
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleWaterLog(context.Background(), nil, testUpdateFromUser(42, 1001, buttonWaterLog))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	if !strings.Contains(gotText, "Журнал полива:") {
		t.Fatalf("expected water log title, got %q", gotText)
	}

	if !strings.Contains(gotText, "Monstera") {
		t.Fatalf("expected plant name in text, got %q", gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, careMenuKeyboard()) {
		t.Fatalf("expected care menu keyboard, got %#v", gotKeyboard)
	}
}

func testUpdate(chatID int64) *models.Update {
	return &models.Update{
		Message: &models.Message{
			Chat: models.Chat{
				ID: chatID,
			},
		},
	}
}

func TestHandlePlantsSendTextWithKeyboard(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handlePlants(context.Background(), nil, testUpdate(42))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	wantText := "Раздел растений 🌱\n\nВыбери действие."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, plantsMenuKeyboard()) {
		t.Fatalf("expected plants menu keyboard, got %#v", gotKeyboard)
	}
}

func TestHandleAddPlantsSendTextWithKeyboard(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup

	b := &Bot{
		log:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		states: NewStateStore(),
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleAddPlant(context.Background(), nil, testUpdateFromUser(42, 1001, buttonAddPlant))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	wantText := "Введи имя растения.\n\nЧтобы выйти, нажми «Отмена» или «Меню»."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, cancelKeyboard()) {
		t.Fatalf("expected cancel menu keyboard, got %#v", gotKeyboard)
	}

	if newState := b.states.Get(1001); newState != StateWaitingPlantName {
		t.Fatalf("expected state %q, got %q", StateWaitingPlantName, newState)
	}
}

func TestHandleListPlants(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	var gotTelegramUserID int64
	var gotUserID int64

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				gotTelegramUserID = telegramUserID
				return domain.User{
					ID:             77,
					TelegramUserID: telegramUserID,
				}, nil
			},
		},
		plants: plantUsecaseStub{
			listPlantsFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
				gotUserID = userID
				return []domain.Plant{
					{ID: 1, UserID: userID, Name: "Monstera"},
					{ID: 2, UserID: userID, Name: "Cactus"},
				}, nil
			},
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleListPlants(context.Background(), nil, testUpdateFromUser(42, 1001, buttonListPlants))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	if gotTelegramUserID != 1001 {
		t.Fatalf("expected EnsureUser telegram user ID %d, got %d", 1001, gotTelegramUserID)
	}

	if gotUserID != 77 {
		t.Fatalf("expected ListPlants user ID %d, got %d", 77, gotUserID)
	}

	wantText := "Твои растения:\n1. Monstera\n2. Cactus"
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, plantsMenuKeyboard()) {
		t.Fatalf("expected plants menu keyboard, got %#v", gotKeyboard)
	}
}

func TestHandleListPlantsEnsureUserError(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	listPlantsCalled := false

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{}, domain.ValidationError{
					Field:   "telegramUserID",
					Problem: "must be positive",
				}
			},
		},
		plants: plantUsecaseStub{
			listPlantsFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
				listPlantsCalled = true
				return nil, nil
			},
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleListPlants(context.Background(), nil, testUpdateFromUser(42, 1001, buttonListPlants))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	wantText := "Не удалось определить пользователя. Попробуй ещё раз позже."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, plantsMenuKeyboard()) {
		t.Fatalf("expected plants menu keyboard, got %#v", gotKeyboard)
	}

	if listPlantsCalled {
		t.Fatal("ListPlants should not be called when EnsureUser fails")
	}
}

func TestHandleDeletePlantStartsDeleteFlow(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.InlineKeyboardMarkup

	b := &Bot{
		log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		pendingDeletes: NewPendingDeleteStore(),
		states:         NewStateStore(),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{ID: 77, TelegramUserID: telegramUserID}, nil
			},
		},
		plants: plantUsecaseStub{
			listPlantsFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
				return []domain.Plant{
					{ID: 1, UserID: userID, Name: "Monstera"},
					{ID: 2, UserID: userID, Name: "Cactus"},
				}, nil
			},
		},
		sendTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.InlineKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	b.handleDeletePlant(context.Background(), nil, testUpdateFromUser(42, 1001, buttonDeletePlant))

	if gotChatID != 42 {
		t.Fatalf("expected chat ID %d, got %d", 42, gotChatID)
	}

	wantText := "Выбери растение для удаления:"
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if gotKeyboard.InlineKeyboard[0][0].CallbackData != callbackDeleteSelectPrefix+"1" {
		t.Fatalf("unexpected first delete button: %#v", gotKeyboard.InlineKeyboard[0][0])
	}
}

func TestHandleDeletePlantSelectionAndConfirm(t *testing.T) {
	var gotEditChatID int64
	var gotEditMessageID int
	var gotText string
	var gotKeyboard models.InlineKeyboardMarkup
	var gotAnsweredCallbackID string

	b := &Bot{
		log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		pendingDeletes: NewPendingDeleteStore(),
		states:         NewStateStore(),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{ID: 77, TelegramUserID: telegramUserID}, nil
			},
		},
		plants: plantUsecaseStub{
			getPlantFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
				return domain.Plant{ID: plantID, UserID: userID, Name: "Monstera"}, nil
			},
			deletePlantFn: func(ctx context.Context, userID int64, plantID int64) error {
				if userID != 77 || plantID != 1 {
					t.Fatalf("unexpected delete params userID=%d plantID=%d", userID, plantID)
				}
				return nil
			},
		},
		editMessageTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
			gotEditChatID = chatID
			gotEditMessageID = messageID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
		answerCallbackQueryFn: func(_ context.Context, callbackQueryID string) error {
			gotAnsweredCallbackID = callbackQueryID
			return nil
		},
	}

	b.handleDeleteSelectCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-1", callbackDeleteSelectPrefix+"1"))

	if gotEditChatID != 42 || gotEditMessageID != 55 {
		t.Fatalf("expected edited message 42/55, got %d/%d", gotEditChatID, gotEditMessageID)
	}

	if gotAnsweredCallbackID != "cb-1" {
		t.Fatalf("expected answered callback %q, got %q", "cb-1", gotAnsweredCallbackID)
	}

	if gotText != `Удалить растение "Monstera"?` {
		t.Fatalf("expected confirm text, got %q", gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, deleteConfirmInlineKeyboard(1)) {
		t.Fatalf("expected delete confirm keyboard, got %#v", gotKeyboard)
	}

	pending, ok := b.pendingDeleteStore().Get(1001, 55)
	if !ok || pending.plantID != 1 {
		t.Fatalf("expected pending delete for plant 1, got %+v, ok=%v", pending, ok)
	}

	b.handleDeleteConfirmCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-2", callbackDeleteConfirmPrefix+"1"))

	if gotText != `Растение "Monstera" удалено.` {
		t.Fatalf("expected delete confirmation text, got %q", gotText)
	}

	if len(gotKeyboard.InlineKeyboard) != 0 {
		t.Fatalf("expected empty inline keyboard after delete, got %#v", gotKeyboard)
	}
}

func TestHandleDeleteCancelCallback(t *testing.T) {
	var gotText string
	var gotKeyboard models.InlineKeyboardMarkup

	b := &Bot{
		log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		pendingDeletes: NewPendingDeleteStore(),
		editMessageTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
		answerCallbackQueryFn: func(_ context.Context, callbackQueryID string) error {
			return nil
		},
	}

	b.pendingDeleteStore().Set(1001, 55, pendingDelete{userID: 77, plantID: 1, plantName: "Monstera"})
	b.handleDeleteCancelCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-cancel", callbackDeleteCancel))

	if gotText != "Удаление отменено." {
		t.Fatalf("expected cancel text, got %q", gotText)
	}

	if len(gotKeyboard.InlineKeyboard) != 0 {
		t.Fatalf("expected empty inline keyboard after cancel, got %#v", gotKeyboard)
	}
}

func TestHandleCareSelectWaterAndBackCallbacks(t *testing.T) {
	var gotEditChatID int64
	var gotEditMessageID int
	var gotText string
	var gotKeyboard models.InlineKeyboardMarkup
	var answeredCallbacks []string

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{ID: 77, TelegramUserID: telegramUserID}, nil
			},
		},
		plants: plantUsecaseStub{
			getPlantFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
				return domain.Plant{ID: plantID, UserID: userID, Name: "Monstera"}, nil
			},
			listPlantsFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
				return []domain.Plant{
					{ID: 1, UserID: userID, Name: "Monstera"},
					{ID: 2, UserID: userID, Name: "Cactus"},
				}, nil
			},
		},
		care: careUsecaseStub{
			markWaterFn: func(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
				if userID != 77 || plantID != 1 {
					t.Fatalf("unexpected MarkWater params userID=%d plantID=%d", userID, plantID)
				}

				return domain.CareEvent{ID: 5, PlantID: plantID, Kind: domain.CareKindWater}, nil
			},
		},
		editMessageTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
			gotEditChatID = chatID
			gotEditMessageID = messageID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
		answerCallbackQueryFn: func(_ context.Context, callbackQueryID string) error {
			answeredCallbacks = append(answeredCallbacks, callbackQueryID)
			return nil
		},
	}

	b.handleCareSelectCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-care-select", callbackCareSelectPrefix+"1"))

	if gotEditChatID != 42 || gotEditMessageID != 55 {
		t.Fatalf("expected edited message 42/55, got %d/%d", gotEditChatID, gotEditMessageID)
	}

	if gotText != `Что отметить для "Monstera"?` {
		t.Fatalf("expected care actions prompt, got %q", gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, careActionsInlineKeyboard(1)) {
		t.Fatalf("expected care action keyboard, got %#v", gotKeyboard)
	}

	b.handleCareWaterCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-care-water", callbackCareWaterPrefix+"1"))

	if gotText != "Полив для Monstera отмечен." {
		t.Fatalf("expected water confirmation, got %q", gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, careActionsInlineKeyboard(1)) {
		t.Fatalf("expected care action keyboard after water, got %#v", gotKeyboard)
	}

	b.handleCareBackCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-care-back", callbackCareBack))

	if gotText != formatCarePlantPrompt() {
		t.Fatalf("expected care plants prompt after back, got %q", gotText)
	}

	if gotKeyboard.InlineKeyboard[0][0].CallbackData != callbackCareSelectPrefix+"1" {
		t.Fatalf("unexpected first care button after back: %#v", gotKeyboard.InlineKeyboard[0][0])
	}

	if len(answeredCallbacks) != 3 {
		t.Fatalf("expected 3 answered callbacks, got %d", len(answeredCallbacks))
	}
}

func TestHandleDeleteConfirmCallbackUsesMessageScopedPendingDelete(t *testing.T) {
	var gotText string

	b := &Bot{
		log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		pendingDeletes: NewPendingDeleteStore(),
		editMessageTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
			gotText = text
			return nil
		},
		answerCallbackQueryFn: func(_ context.Context, callbackQueryID string) error {
			return nil
		},
	}

	b.pendingDeleteStore().Set(1001, 99, pendingDelete{userID: 77, plantID: 1, plantName: "Monstera"})

	b.handleDeleteConfirmCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-2", callbackDeleteConfirmPrefix+"1"))

	if gotText != "Не удалось продолжить удаление. Попробуй заново." {
		t.Fatalf("expected missing pending delete text, got %q", gotText)
	}
}

func TestHandleDeleteConfirmCallbackEditsMessageOnDeleteError(t *testing.T) {
	var gotText string
	var sendReplyCalled bool

	b := &Bot{
		log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		pendingDeletes: NewPendingDeleteStore(),
		plants: plantUsecaseStub{
			deletePlantFn: func(ctx context.Context, userID int64, plantID int64) error {
				return domain.ErrNotFound
			},
		},
		editMessageTextWithInlineKeyboardFn: func(_ context.Context, chatID int64, messageID int, text string, keyboard models.InlineKeyboardMarkup) error {
			gotText = text
			return nil
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			sendReplyCalled = true
			return nil
		},
		answerCallbackQueryFn: func(_ context.Context, callbackQueryID string) error {
			return nil
		},
	}

	b.pendingDeleteStore().Set(1001, 55, pendingDelete{userID: 77, plantID: 1, plantName: "Monstera"})

	b.handleDeleteConfirmCallback(context.Background(), nil, testDeleteCallbackUpdate(42, 1001, 55, "cb-2", callbackDeleteConfirmPrefix+"1"))

	if gotText != "Растение не найдено." {
		t.Fatalf("expected edited callback error text, got %q", gotText)
	}

	if sendReplyCalled {
		t.Fatal("callback delete error should not send a new reply-keyboard message")
	}
}

func TestHandleAddPlantValidName(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	var gotTelegramUserID int64
	var gotUserID int64
	var gotName string
	ctx := context.Background()

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				gotTelegramUserID = telegramUserID
				return domain.User{
					ID:             77,
					TelegramUserID: telegramUserID,
					CreatedAt:      time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC),
				}, nil
			},
		},
		plants: plantUsecaseStub{
			addPlantFn: func(ctx context.Context, userID int64, name string) (domain.Plant, error) {
				gotUserID = userID
				gotName = name
				return domain.Plant{
					ID:     1,
					UserID: userID,
					Name:   "Фикус",
				}, nil
			},
		},
		states: NewStateStore(),
		sendTextFn: func(_ context.Context, chatID int64, text string) error {
			gotChatID = chatID
			gotText = text
			return nil
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	const (
		chatID = int64(42)
		userID = int64(1001)
	)

	b.handleAddPlant(ctx, nil, testUpdateFromUser(chatID, userID, buttonAddPlant))
	if newState := b.states.Get(userID); newState != StateWaitingPlantName {
		t.Fatalf("expected state %q, got %q", StateWaitingPlantName, newState)
	}

	b.handleTextByState(ctx, nil, testUpdateFromUser(chatID, userID, "Фикус"))
	if gotChatID != chatID {
		t.Fatalf("expected chat ID %d, got %d", chatID, gotChatID)
	}

	if gotTelegramUserID != userID {
		t.Fatalf("expected EnsureUser telegram user ID %d, got %d", userID, gotTelegramUserID)
	}

	if gotUserID != 77 {
		t.Fatalf("expected AddPlant user ID %d, got %d", 77, gotUserID)
	}

	if gotName != "Фикус" {
		t.Fatalf("expected AddPlant name %q, got %q", "Фикус", gotName)
	}

	wantText := `Растение "Фикус" добавлено 🌿`
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, plantsMenuKeyboard()) {
		t.Fatalf("expected plants keyboard, got %#v", gotKeyboard)
	}

	if b.states.Get(userID) != StateIdle {
		t.Fatalf("expected state %q, got %q", StateIdle, b.states.Get(userID))
	}
}

func TestHandleAddPlantEmptyName(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	ctx := context.Background()

	b := &Bot{
		log:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		states: NewStateStore(),
		sendTextFn: func(_ context.Context, chatID int64, text string) error {
			gotChatID = chatID
			gotText = text
			return nil
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	const (
		chatID = int64(42)
		userID = int64(1001)
	)

	b.handleAddPlant(ctx, nil, testUpdateFromUser(chatID, userID, buttonAddPlant))
	b.handlePlantNameInput(ctx, chatID, userID, "    ")

	if gotChatID != chatID {
		t.Fatalf("expected chat ID %d, got %d", chatID, gotChatID)
	}

	wantText := "Имя растения не должно быть пустым."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, cancelKeyboard()) {
		t.Fatalf("expected cancel keyboard, got %#v", gotKeyboard)
	}

	if newState := b.states.Get(userID); newState != StateWaitingPlantName {
		t.Fatalf("expected state %q, got %q", StateWaitingPlantName, newState)
	}
}

func TestHandleAddPlantEnsureUserError(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	addPlantCalled := false
	ctx := context.Background()

	b := &Bot{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		users: userUsecaseStub{
			ensureUserFn: func(ctx context.Context, telegramUserID int64) (domain.User, error) {
				return domain.User{}, context.DeadlineExceeded
			},
		},
		plants: plantUsecaseStub{
			addPlantFn: func(ctx context.Context, userID int64, name string) (domain.Plant, error) {
				addPlantCalled = true
				return domain.Plant{}, nil
			},
		},
		states: NewStateStore(),
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	const (
		chatID = int64(42)
		userID = int64(1001)
	)

	b.handleAddPlant(ctx, nil, testUpdateFromUser(chatID, userID, buttonAddPlant))
	b.handlePlantNameInput(ctx, chatID, userID, "Фикус")

	if gotChatID != chatID {
		t.Fatalf("expected chat ID %d, got %d", chatID, gotChatID)
	}

	if gotText != "Что-то пошло не так. Попробуй ещё раз позже." {
		t.Fatalf("expected generic error text, got %q", gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, cancelKeyboard()) {
		t.Fatalf("expected cancel keyboard, got %#v", gotKeyboard)
	}

	if addPlantCalled {
		t.Fatal("AddPlant should not be called when EnsureUser fails")
	}

	if newState := b.states.Get(userID); newState != StateWaitingPlantName {
		t.Fatalf("expected state %q, got %q", StateWaitingPlantName, newState)
	}
}

func TestHandleCancelClearsState(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	ctx := context.Background()

	b := &Bot{
		log:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		states: NewStateStore(),
		sendTextFn: func(_ context.Context, chatID int64, text string) error {
			gotChatID = chatID
			gotText = text
			return nil
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	const (
		chatID = int64(42)
		userID = int64(1001)
	)

	b.states.Set(userID, StateWaitingPlantName)
	b.handleCancel(ctx, nil, testUpdateFromUser(chatID, userID, buttonCancel))

	if gotChatID != chatID {
		t.Fatalf("expected chat ID %d, got %d", chatID, gotChatID)
	}

	wantText := "Действие отменено."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, plantsMenuKeyboard()) {
		t.Fatalf("expected plants menu keyboard, got %#v", gotKeyboard)
	}

	if newState := b.states.Get(userID); newState != StateIdle {
		t.Fatalf("expected state %q, got %q", StateIdle, newState)
	}
}

func TestHandleBackClearsState(t *testing.T) {
	var gotChatID int64
	var gotText string
	var gotKeyboard models.ReplyKeyboardMarkup
	ctx := context.Background()

	b := &Bot{
		log:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		states: NewStateStore(),
		sendTextFn: func(_ context.Context, chatID int64, text string) error {
			gotChatID = chatID
			gotText = text
			return nil
		},
		sendTextWithKeyboardFn: func(_ context.Context, chatID int64, text string, keyboard models.ReplyKeyboardMarkup) error {
			gotChatID = chatID
			gotText = text
			gotKeyboard = keyboard
			return nil
		},
	}

	const (
		chatID = int64(42)
		userID = int64(1001)
	)

	b.states.Set(userID, StateWaitingPlantName)
	b.handleBackToMenu(ctx, nil, testUpdateFromUser(chatID, userID, buttonBackToMenu))

	if gotChatID != chatID {
		t.Fatalf("expected chat ID %d, got %d", chatID, gotChatID)
	}

	wantText := "Возвращаю в главное меню."
	if gotText != wantText {
		t.Fatalf("expected text %q, got %q", wantText, gotText)
	}

	if !reflect.DeepEqual(gotKeyboard, mainMenuKeyboard()) {
		t.Fatalf("expected main menu keyboard, got %#v", gotKeyboard)
	}

	if newState := b.states.Get(userID); newState != StateIdle {
		t.Fatalf("expected state %q, got %q", StateIdle, newState)
	}
}

func testUpdateFromUser(chatID, userID int64, text string) *models.Update {
	return &models.Update{
		Message: &models.Message{
			Chat: models.Chat{
				ID: chatID,
			},
			From: &models.User{
				ID: userID,
			},
			Text: text,
		},
	}
}

func testDeleteCallbackUpdate(chatID, userID int64, messageID int, callbackID, data string) *models.Update {
	return &models.Update{
		CallbackQuery: &models.CallbackQuery{
			ID:   callbackID,
			From: models.User{ID: userID},
			Data: data,
			Message: models.MaybeInaccessibleMessage{
				Type: models.MaybeInaccessibleMessageTypeMessage,
				Message: &models.Message{
					ID: messageID,
					Chat: models.Chat{
						ID: chatID,
					},
				},
			},
		},
	}
}
