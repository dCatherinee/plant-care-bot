package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewCareEvent(t *testing.T) {
	loc := time.FixedZone("MSK", 3*60*60)         // UTC+3
	in := time.Date(2026, 3, 4, 12, 0, 0, 0, loc) // 12:00 MSK
	fixed := time.Date(2026, 3, 4, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		plantID     int64
		kind        CareKind
		occurredAt  time.Time
		wantErr     bool
		wantField   string
		wantProblem string
	}{
		{"ok", 10, CareKindWater, fixed, false, "", ""},
		{"empty_plant_id", 0, CareKindWater, fixed, true, "plantID", "must be positive"},
		{"invalid_kind", 10, CareKind("freeze"), fixed, true, "kind", "invalid care kind"},
		{"invalid_location", 10, CareKindFertilize, in, false, "", ""},
		{"zero_time", 10, CareKindWater, time.Time{}, true, "occurredAt", "is zero"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			event, err := NewCareEvent(tc.plantID, tc.kind, tc.occurredAt)

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
			if event.OccurredAt.Location() != time.UTC {
				t.Fatalf("OccurredAt must be UTC, got %v", event.OccurredAt.Location())
			}
			if !event.OccurredAt.Equal(tc.occurredAt.UTC()) {
				t.Fatalf("OccurredAt mismatch: got %v want %v", event.OccurredAt, tc.occurredAt.UTC())
			}
		})
	}
}
