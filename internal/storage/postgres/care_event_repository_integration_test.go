//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func createTestCareEvent(t *testing.T, plantID int64, kind domain.CareKind, occurredAt time.Time) int64 {
	t.Helper()

	repo := NewCareEventRepository(newTestDB(t))
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	event, err := repo.CreateCareEvent(ctx, domain.CareEvent{
		PlantID:    plantID,
		Kind:       kind,
		OccurredAt: occurredAt,
	})
	if err != nil {
		t.Fatalf("create test care event: %v", err)
	}

	return event.ID
}

func TestCareEventRepositoryCreateCareEvent_Integration(t *testing.T) {
	db := newTestDB(t)
	cleanupTables(t, db)

	userID := createTestUser(t, 4001)
	plantID := createTestPlant(t, userID, "Monstera", time.Date(2026, time.April, 17, 9, 0, 0, 0, time.UTC))
	repo := NewCareEventRepository(db)
	ctx := context.Background()
	occurredAt := time.Date(2026, time.April, 17, 10, 0, 0, 0, time.UTC)

	event := domain.CareEvent{
		PlantID:    plantID,
		Kind:       domain.CareKindWater,
		OccurredAt: occurredAt,
	}

	savedEvent, err := repo.CreateCareEvent(ctx, event)
	if err != nil {
		t.Fatalf("CreateCareEvent returned error: %v", err)
	}

	if savedEvent.ID <= 0 {
		t.Fatalf("expected positive event ID, got %d", savedEvent.ID)
	}

	const query = `
		select id, plant_id, event_type, occurred_at, created_at
		from care_events
		where id = $1
	`

	var persisted domain.CareEvent
	if err := db.QueryRowContext(ctx, query, savedEvent.ID).Scan(
		&persisted.ID,
		&persisted.PlantID,
		&persisted.Kind,
		&persisted.OccurredAt,
		&persisted.CreatedAt,
	); err != nil {
		t.Fatalf("query saved care event: %v", err)
	}

	if persisted.ID != savedEvent.ID {
		t.Fatalf("expected saved id %d, got %d", savedEvent.ID, persisted.ID)
	}
	if persisted.PlantID != event.PlantID {
		t.Fatalf("expected plant id %d, got %d", event.PlantID, persisted.PlantID)
	}
	if persisted.Kind != event.Kind {
		t.Fatalf("expected kind %q, got %q", event.Kind, persisted.Kind)
	}
	if !persisted.OccurredAt.Equal(event.OccurredAt) {
		t.Fatalf("expected occurred_at %v, got %v", event.OccurredAt, persisted.OccurredAt)
	}
}

func TestCareEventRepositoryListCareEventsByType_Integration(t *testing.T) {
	db := newTestDB(t)
	cleanupTables(t, db)

	userID := createTestUser(t, 4002)
	plantID := createTestPlant(t, userID, "Cactus", time.Date(2026, time.April, 17, 9, 0, 0, 0, time.UTC))
	firstOccurredAt := time.Date(2026, time.April, 17, 12, 0, 0, 0, time.UTC)
	secondOccurredAt := time.Date(2026, time.April, 17, 11, 0, 0, 0, time.UTC)
	createTestCareEvent(t, plantID, domain.CareKindWater, firstOccurredAt)
	createTestCareEvent(t, plantID, domain.CareKindWater, secondOccurredAt)
	createTestCareEvent(t, plantID, domain.CareKindFertilize, time.Date(2026, time.April, 17, 10, 0, 0, 0, time.UTC))

	repo := NewCareEventRepository(db)
	ctx := context.Background()

	events, err := repo.ListCareEventsByType(ctx, plantID, domain.CareKindWater)
	if err != nil {
		t.Fatalf("ListCareEventsByType returned error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].Kind != domain.CareKindWater || events[1].Kind != domain.CareKindWater {
		t.Fatalf("expected only water events, got %+v", events)
	}

	if !events[0].OccurredAt.Equal(firstOccurredAt) {
		t.Fatalf("expected first occurred_at %v, got %v", firstOccurredAt, events[0].OccurredAt)
	}
	if !events[1].OccurredAt.Equal(secondOccurredAt) {
		t.Fatalf("expected second occurred_at %v, got %v", secondOccurredAt, events[1].OccurredAt)
	}
}

func TestCareEventRepositoryListRecentCareEventsByUserAndType_Integration(t *testing.T) {
	db := newTestDB(t)
	cleanupTables(t, db)

	firstUserID := createTestUser(t, 4003)
	secondUserID := createTestUser(t, 4004)
	firstPlantID := createTestPlant(t, firstUserID, "Monstera", time.Date(2026, time.April, 17, 8, 0, 0, 0, time.UTC))
	secondPlantID := createTestPlant(t, firstUserID, "Cactus", time.Date(2026, time.April, 17, 8, 30, 0, 0, time.UTC))
	otherUserPlantID := createTestPlant(t, secondUserID, "Orchid", time.Date(2026, time.April, 17, 9, 0, 0, 0, time.UTC))

	latestOccurredAt := time.Date(2026, time.April, 17, 15, 0, 0, 0, time.UTC)
	middleOccurredAt := time.Date(2026, time.April, 17, 14, 0, 0, 0, time.UTC)
	oldestOccurredAt := time.Date(2026, time.April, 17, 13, 0, 0, 0, time.UTC)

	createTestCareEvent(t, firstPlantID, domain.CareKindWater, oldestOccurredAt)
	createTestCareEvent(t, secondPlantID, domain.CareKindWater, latestOccurredAt)
	createTestCareEvent(t, firstPlantID, domain.CareKindWater, middleOccurredAt)
	createTestCareEvent(t, firstPlantID, domain.CareKindFertilize, time.Date(2026, time.April, 17, 12, 0, 0, 0, time.UTC))
	createTestCareEvent(t, otherUserPlantID, domain.CareKindWater, time.Date(2026, time.April, 17, 16, 0, 0, 0, time.UTC))

	repo := NewCareEventRepository(db)
	ctx := context.Background()

	events, err := repo.ListRecentCareEventsByUserAndType(ctx, firstUserID, domain.CareKindWater, 2)
	if err != nil {
		t.Fatalf("ListRecentCareEventsByUserAndType returned error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if !events[0].OccurredAt.Equal(latestOccurredAt) {
		t.Fatalf("expected latest occurred_at %v, got %v", latestOccurredAt, events[0].OccurredAt)
	}
	if !events[1].OccurredAt.Equal(middleOccurredAt) {
		t.Fatalf("expected second occurred_at %v, got %v", middleOccurredAt, events[1].OccurredAt)
	}

	for _, event := range events {
		if event.Kind != domain.CareKindWater {
			t.Fatalf("expected only water events, got %+v", events)
		}
		if event.PlantID != firstPlantID && event.PlantID != secondPlantID {
			t.Fatalf("expected only first user plants, got plantID=%d", event.PlantID)
		}
	}
}
