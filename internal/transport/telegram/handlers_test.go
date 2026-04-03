package telegram

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"testing"

	"github.com/go-telegram/bot/models"
)

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
			name: "care",
			handler: func(ctx context.Context, bt *Bot, update *models.Update) {
				bt.handleCare(ctx, nil, update)
			},
			want: `Раздел "Уход" пока в разработке 💧`,
		},
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

func TestHandleAddPlantValidName(t *testing.T) {
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
	if newState := b.states.Get(userID); newState != StateWaitingPlantName {
		t.Fatalf("expected state %q, got %q", StateWaitingPlantName, newState)
	}

	b.handleTextByState(ctx, nil, testUpdateFromUser(chatID, userID, "Фикус"))
	if gotChatID != chatID {
		t.Fatalf("expected chat ID %d, got %d", chatID, gotChatID)
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

	wantText := "Имя растения не должно быть пустым. Введи название ещё раз."
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
