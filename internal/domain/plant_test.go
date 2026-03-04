package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestPlant_EmptyUserID(t *testing.T) {
	_, err := NewPlant(0, "cactus")

	if !errors.Is(err, ErrInvalidArgument) {
		t.Error("UserID can't be empty or negative number")
	}
}

func TestPlant_EmptyName(t *testing.T) {
	_, err := NewPlant(10, "")

	if !errors.Is(err, ErrInvalidArgument) {
		t.Error("Plant name can't be empty")
	}
}

func TestPlant_TrimName(t *testing.T) {
	plantName := " Cactus Poppy  "

	plant, err := NewPlant(10, plantName)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plant.Name != strings.TrimSpace(plantName) {
		t.Error("Plant name don't trim")
	}
}

func TestPlant_InvalidCreateDate(t *testing.T) {
	plant, err := NewPlant(10, "cactus")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plant.CreatedAt.IsZero() {
		t.Error("Plant created date can't be empty")
	}

	if plant.CreatedAt.Location() != time.UTC {
		t.Fatalf("CreatedAt must be UTC, got %v", plant.CreatedAt.Location())
	}
}
