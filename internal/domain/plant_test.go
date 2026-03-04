package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestPlant_EmptyUserID(t *testing.T) {
	_, err := NewPlant(0, "cactus")

	if !errors.Is(err, ErrInvalidArgument) {
		t.Error("UserID can't be empty or negative number")
	}
}

func TestPlant_EmptyName(t *testing.T) {
	_, err := NewPlant(10, "")

	if err == nil {
		t.Error("Plant name can't be empty")
	}
}

func TestPlant_TrimName(t *testing.T) {
	plantName := " Cactus Poppy  "

	plant, _ := NewPlant(10, plantName)

	if plant.Name != strings.TrimSpace(plantName) {
		t.Error("Plant name don't trim")
	}
}

func TestPlant_InvalidCreateDate(t *testing.T) {
	plant, _ := NewPlant(10, "cactus")

	if plant.CreatedAt.IsZero() {
		t.Error("Plant created date can't be empty")
	}
}
