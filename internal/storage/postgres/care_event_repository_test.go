package postgres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func newCareEventRepositoryTestDB(t *testing.T) (*CareEventRepository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("db.Close: %v", err)
		}
	})

	return NewCareEventRepository(db), mock
}

func TestCareEventRepositoryCreateCareEventReturnsEvent(t *testing.T) {
	repo, mock := newCareEventRepositoryTestDB(t)
	ctx := context.Background()
	occurredAt := time.Date(2026, time.April, 17, 10, 0, 0, 0, time.UTC)
	createdAt := occurredAt.Add(5 * time.Minute)
	input := domain.CareEvent{
		PlantID:    15,
		Kind:       domain.CareKindWater,
		OccurredAt: occurredAt,
	}
	expected := domain.CareEvent{
		ID:         101,
		PlantID:    15,
		Kind:       domain.CareKindWater,
		OccurredAt: occurredAt,
		CreatedAt:  createdAt,
	}

	query := regexp.QuoteMeta(`
		insert into care_events (plant_id, event_type, occurred_at)
		values ($1, $2, $3)
		returning id, plant_id, event_type, occurred_at, created_at;
	`)

	mock.ExpectQuery(query).
		WithArgs(input.PlantID, input.Kind, input.OccurredAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "plant_id", "event_type", "occurred_at", "created_at"}).
			AddRow(expected.ID, expected.PlantID, expected.Kind, expected.OccurredAt, expected.CreatedAt))

	actual, err := repo.CreateCareEvent(ctx, input)
	if err != nil {
		t.Fatalf("CreateCareEvent returned error: %v", err)
	}

	if actual != expected {
		t.Fatalf("expected care event %+v, got %+v", expected, actual)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestCareEventRepositoryListCareEventsByTypeMapsRows(t *testing.T) {
	repo, mock := newCareEventRepositoryTestDB(t)
	ctx := context.Background()
	plantID := int64(15)
	firstOccurredAt := time.Date(2026, time.April, 17, 11, 0, 0, 0, time.UTC)
	secondOccurredAt := time.Date(2026, time.April, 17, 10, 0, 0, 0, time.UTC)
	firstCreatedAt := firstOccurredAt.Add(5 * time.Minute)
	secondCreatedAt := secondOccurredAt.Add(5 * time.Minute)

	query := regexp.QuoteMeta(`
		select id, plant_id, event_type, occurred_at, created_at
		from care_events
		where plant_id = $1 and event_type = $2
		order by occurred_at desc
	`)

	rows := sqlmock.NewRows([]string{"id", "plant_id", "event_type", "occurred_at", "created_at"}).
		AddRow(int64(2), plantID, domain.CareKindWater, firstOccurredAt, firstCreatedAt).
		AddRow(int64(1), plantID, domain.CareKindWater, secondOccurredAt, secondCreatedAt)

	mock.ExpectQuery(query).
		WithArgs(plantID, domain.CareKindWater).
		WillReturnRows(rows)

	events, err := repo.ListCareEventsByType(ctx, plantID, domain.CareKindWater)
	if err != nil {
		t.Fatalf("ListCareEventsByType returned error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].ID != 2 || events[1].ID != 1 {
		t.Fatalf("unexpected events returned: %+v", events)
	}

	if events[0].Kind != domain.CareKindWater || events[1].Kind != domain.CareKindWater {
		t.Fatalf("unexpected care kind in result: %+v", events)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestCareEventRepositoryListRecentCareEventsByUserAndTypeMapsRows(t *testing.T) {
	repo, mock := newCareEventRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(77)
	limit := 3
	firstOccurredAt := time.Date(2026, time.April, 17, 11, 0, 0, 0, time.UTC)
	secondOccurredAt := time.Date(2026, time.April, 17, 10, 0, 0, 0, time.UTC)
	firstCreatedAt := firstOccurredAt.Add(5 * time.Minute)
	secondCreatedAt := secondOccurredAt.Add(5 * time.Minute)

	query := regexp.QuoteMeta(`
		select ce.id, ce.plant_id, ce.event_type, ce.occurred_at, ce.created_at
		from care_events ce
		join plants p on p.id = ce.plant_id
		where p.user_id = $1 and ce.event_type = $2
		order by ce.occurred_at desc
		limit $3
	`)

	rows := sqlmock.NewRows([]string{"id", "plant_id", "event_type", "occurred_at", "created_at"}).
		AddRow(int64(3), int64(2), domain.CareKindFertilize, firstOccurredAt, firstCreatedAt).
		AddRow(int64(1), int64(1), domain.CareKindFertilize, secondOccurredAt, secondCreatedAt)

	mock.ExpectQuery(query).
		WithArgs(userID, domain.CareKindFertilize, limit).
		WillReturnRows(rows)

	events, err := repo.ListRecentCareEventsByUserAndType(ctx, userID, domain.CareKindFertilize, limit)
	if err != nil {
		t.Fatalf("ListRecentCareEventsByUserAndType returned error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].ID != 3 || events[1].ID != 1 {
		t.Fatalf("unexpected events returned: %+v", events)
	}

	if events[0].Kind != domain.CareKindFertilize || events[1].Kind != domain.CareKindFertilize {
		t.Fatalf("unexpected care kind in result: %+v", events)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
