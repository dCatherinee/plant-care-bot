package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewPlant(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		plantName string
		wantErr   bool
	}{
		{"ok", 1, "Monstera", false},
		{"empty_user_id", 0, "Cactus", true},
		{"empty_name", 10, "", true},
		{"trim_name", 10, " Cactus Poppy  ", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plant, err := NewPlant(tc.userID, tc.plantName)

			if tc.wantErr {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}

				var myErr ValidationError
				if errors.As(err, &myErr) {
					if myErr.Field == "name" {
						t.Fatalf("Expected error on field: %q: %q", myErr.Field, myErr.Problem)
					}
				}
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("Expected ErrInvalidArgument, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if strings.TrimSpace(tc.plantName) != plant.Name {
				t.Fatalf("Expected trimmed name %q, got %q", strings.TrimSpace(tc.plantName), plant.Name)
			}
			if plant.CreatedAt.Location() != time.UTC {
				t.Fatalf("CreatedAt must be UTC, got %v", plant.CreatedAt.Location())
			}
			if plant.CreatedAt.IsZero() {
				t.Fatal("CreatedAt must not be zero")
			}
		})
	}
}
