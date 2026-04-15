package telegram

import "testing"

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

	if len(keyboard.Keyboard) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(keyboard.Keyboard))
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
