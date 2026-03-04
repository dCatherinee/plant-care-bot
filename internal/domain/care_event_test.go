package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCareEvent_EmptyPlantID(t *testing.T) {
	_, err := NewCareEvent(0, CareKindWater, time.Now())

	if !errors.Is(err, ErrInvalidArgument) {
		t.Error("Plant ID can't be empty or negative number")
	}
}

func TestCareEvent_InvalidKind(t *testing.T) {
	_, err := NewCareEvent(10, "freeze", time.Now())

	if !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("Plant kind can't be %q, allowed only 'water' and 'fertilize'", "freeze")
	}
}

func TestCareEvent_InvalidTime(t *testing.T) {
	loc := time.FixedZone("MSK", 3*60*60)         // UTC+3
	in := time.Date(2026, 3, 4, 12, 0, 0, 0, loc) // 12:00 MSK

	ev, err := NewCareEvent(10, CareKindFertilize, in)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ev.OccurredAt.Location() != time.UTC {
		t.Fatalf("OccurredAt must be UTC, got %v", ev.OccurredAt.Location())
	}
}
