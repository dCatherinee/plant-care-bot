package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCareEvent(t *testing.T) {
	loc := time.FixedZone("MSK", 3*60*60)         // UTC+3
	in := time.Date(2026, 3, 4, 12, 0, 0, 0, loc) // 12:00 MSK
	fixed := time.Date(2026, 3, 4, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		plantID    int64
		kind       CareKind
		occurredAt time.Time
		wantErr    bool
	}{
		{"ok", 10, CareKindWater, fixed, false},
		{"empty_plant_id", 0, CareKindWater, fixed, true},
		{"invalid_kind", 10, CareKind("freeze"), fixed, true},
		{"invalid_location", 10, CareKindFertilize, in, false},
		{"zero_time", 10, CareKindWater, time.Time{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			event, err := NewCareEvent(tc.plantID, tc.kind, tc.occurredAt)

			if tc.wantErr {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidArgument) {
					t.Fatalf("Expected ErrInvalidArgument, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if event.OccurredAt.Location() != time.UTC {
				t.Fatalf("OccurredAt must be UTC, got %v", event.OccurredAt.Location())
			}
			if !event.OccurredAt.Equal(tc.occurredAt.UTC()) {
				t.Fatalf("OccurredAt mismatch: got %v want %v", event.OccurredAt, tc.occurredAt.UTC())
			}
		})
	}
}
