package domain

import (
	"testing"
	"time"
)

func TestCareEvent_EmptyPlantID(t *testing.T) {
	_, err := NewCareEvent(0, "water", time.Now())

	if err == nil {
		t.Error("Plant ID can't be empty or negative number")
	}
}

func TestCareEvent_InvalidKind(t *testing.T) {
	_, err := NewCareEvent(10, "freeze", time.Now())

	if err == nil {
		t.Errorf("Plant kind can't be %q, allowed only 'water' and 'fertilize'", "freeze")
	}
}

func TestCareEvent_InvalidTime(t *testing.T) {
	var occurredAt time.Time
	_, err := NewCareEvent(10, "fertilize", occurredAt)

	if err == nil {
		t.Errorf("Occurred time can't be empty")
	}
}
