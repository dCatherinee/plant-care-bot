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
			name: "plants",
			handler: func(ctx context.Context, bt *Bot, update *models.Update) {
				bt.handlePlants(ctx, nil, update)
			},
			want: `Раздел "Растения" пока в разработке 🌱`,
		},
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
			want: `Раздел "Help" пока в разработке ℹ️`,
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
