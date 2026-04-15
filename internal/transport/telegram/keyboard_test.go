package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestMainMenuKeyboard(t *testing.T) {
	keyboard := mainMenuKeyboard()

	if !keyboard.ResizeKeyboard {
		t.Fatal("expected keyboard to use resize mode")
	}

	if len(keyboard.Keyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(keyboard.Keyboard))
	}

	tests := []struct {
		name string
		row  int
		want []string
	}{
		{
			name: "first row",
			row:  0,
			want: []string{buttonPlants, buttonCare},
		},
		{
			name: "second row",
			row:  1,
			want: []string{buttonReminders, buttonSettings},
		},
		{
			name: "third row",
			row:  2,
			want: []string{buttonHelp},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := keyboard.Keyboard[tt.row]
			if len(row) != len(tt.want) {
				t.Fatalf("expected %d buttons, got %d", len(tt.want), len(row))
			}

			for i, want := range tt.want {
				if row[i].Text != want {
					t.Fatalf("expected button %d text %q, got %q", i, want, row[i].Text)
				}
			}
		})
	}
}

func TestPlantsMenuKeyboard(t *testing.T) {
	keyboard := plantsMenuKeyboard()

	if !keyboard.ResizeKeyboard {
		t.Fatal("expected keyboard to use resize mode")
	}

	if len(keyboard.Keyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(keyboard.Keyboard))
	}

	tests := []struct {
		name string
		row  int
		want []string
	}{
		{
			name: "first row",
			row:  0,
			want: []string{buttonAddPlant, buttonListPlants},
		},
		{
			name: "second row",
			row:  1,
			want: []string{buttonDeletePlant},
		},
		{
			name: "third row",
			row:  2,
			want: []string{buttonBackToMenu},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := keyboard.Keyboard[tt.row]
			if len(row) != len(tt.want) {
				t.Fatalf("expected %d buttons, got %d", len(tt.want), len(row))
			}

			for i, want := range tt.want {
				if row[i].Text != want {
					t.Fatalf("expected button %d text %q, got %q", i, want, row[i].Text)
				}
			}
		})
	}
}

func TestDeletePlantsInlineKeyboard(t *testing.T) {
	keyboard := deletePlantsInlineKeyboard([]models.InlineKeyboardButton{
		{Text: "Monstera", CallbackData: callbackDeleteSelectPrefix + "1"},
		{Text: "Cactus", CallbackData: callbackDeleteSelectPrefix + "2"},
	})

	if len(keyboard.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(keyboard.InlineKeyboard))
	}

	if keyboard.InlineKeyboard[0][0].Text != "Monstera" {
		t.Fatalf("unexpected first button: %#v", keyboard.InlineKeyboard[0][0])
	}

	if keyboard.InlineKeyboard[1][0].Text != "Cactus" {
		t.Fatalf("unexpected second button: %#v", keyboard.InlineKeyboard[1][0])
	}

	if keyboard.InlineKeyboard[2][0].CallbackData != callbackDeleteCancel {
		t.Fatalf("unexpected cancel callback: %#v", keyboard.InlineKeyboard[2][0])
	}
}

func TestDeleteConfirmInlineKeyboard(t *testing.T) {
	keyboard := deleteConfirmInlineKeyboard(12)

	if len(keyboard.InlineKeyboard) != 1 {
		t.Fatalf("expected 1 row, got %d", len(keyboard.InlineKeyboard))
	}

	row := keyboard.InlineKeyboard[0]
	if len(row) != 2 {
		t.Fatalf("expected 2 buttons in row, got %d", len(row))
	}

	if row[0].CallbackData != callbackDeleteConfirmPrefix+"12" {
		t.Fatalf("unexpected confirm callback: %#v", row[0])
	}

	if row[1].CallbackData != callbackDeleteCancel {
		t.Fatalf("unexpected cancel callback: %#v", row[1])
	}
}
