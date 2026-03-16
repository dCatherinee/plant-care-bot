package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewPlant(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		plantName   string
		wantErr     bool
		wantField   string
		wantProblem string
	}{
		{"ok", 1, "Monstera", false, "", ""},
		{"empty_user_id", 0, "Cactus", true, "userID", "must be positive"},
		{"empty_name", 10, "", true, "name", "is empty"},
		{"trim_name", 10, " Cactus Poppy  ", false, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plant, err := NewPlant(tc.userID, tc.plantName)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				var myErr ValidationError
				if !errors.As(err, &myErr) {
					t.Fatalf("expected ValidationError, got %T: %v", err, err)
				}
				if myErr.Field != tc.wantField {
					t.Fatalf("expected field %q, got %q", tc.wantField, myErr.Field)
				}
				if myErr.Problem != tc.wantProblem {
					t.Fatalf("expected problem %q, got %q", tc.wantProblem, myErr.Problem)
				}

				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("expected ErrInvalidArgument, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if strings.TrimSpace(tc.plantName) != plant.Name {
				t.Fatalf("expected trimmed name %q, got %q", strings.TrimSpace(tc.plantName), plant.Name)
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

func TestPlantRename(t *testing.T) {
	tests := []struct {
		name        string
		newName     string
		wantName    string
		wantErr     bool
		wantField   string
		wantProblem string
	}{
		{"ok", "Cactus", "Cactus", false, "", ""},
		{"trim_name", "  Cactus Poppy  ", "Cactus Poppy", false, "", ""},
		{"empty_name", "", "", true, "name", "is empty"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plant := Plant{Name: "Monstera"}

			err := plant.Rename(tc.newName)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				var myErr ValidationError
				if !errors.As(err, &myErr) {
					t.Fatalf("expected ValidationError, got %T: %v", err, err)
				}
				if myErr.Field != tc.wantField {
					t.Fatalf("expected field %q, got %q", tc.wantField, myErr.Field)
				}
				if myErr.Problem != tc.wantProblem {
					t.Fatalf("expected problem %q, got %q", tc.wantProblem, myErr.Problem)
				}
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("expected ErrInvalidArgument, got %v", err)
				}
				if plant.Name != "Monstera" {
					t.Fatalf("expected plant name to stay %q, got %q", "Monstera", plant.Name)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if plant.Name != tc.wantName {
				t.Fatalf("expected renamed plant %q, got %q", tc.wantName, plant.Name)
			}
		})
	}
}
