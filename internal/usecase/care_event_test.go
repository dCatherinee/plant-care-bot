package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type fakeCareEventRepo struct {
	createCareEventFn                 func(ctx context.Context, event domain.CareEvent) (int64, error)
	listCareEventsByTypeFn            func(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error)
	listRecentCareEventsByUserAndType func(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error)
}

func (f *fakeCareEventRepo) CreateCareEvent(ctx context.Context, event domain.CareEvent) (int64, error) {
	if f.createCareEventFn != nil {
		return f.createCareEventFn(ctx, event)
	}

	return 0, nil
}

func (f *fakeCareEventRepo) ListCareEventsByType(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error) {
	if f.listCareEventsByTypeFn != nil {
		return f.listCareEventsByTypeFn(ctx, plantID, eventType)
	}

	return nil, nil
}

func (f *fakeCareEventRepo) ListRecentCareEventsByUserAndType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
	if f.listRecentCareEventsByUserAndType != nil {
		return f.listRecentCareEventsByUserAndType(ctx, userID, eventType, limit)
	}

	return nil, nil
}

func TestCareEventServiceMarkWaterChecksPlantOwnership(t *testing.T) {
	plantRepo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			if userID != 77 || plantID != 15 {
				t.Fatalf("unexpected GetPlantByID params userID=%d plantID=%d", userID, plantID)
			}

			return domain.Plant{ID: plantID, UserID: userID, Name: "Monstera"}, nil
		},
	}

	repo := &fakeCareEventRepo{
		createCareEventFn: func(ctx context.Context, event domain.CareEvent) (int64, error) {
			if event.PlantID != 15 {
				t.Fatalf("expected plant ID 15, got %d", event.PlantID)
			}
			if event.Kind != domain.CareKindWater {
				t.Fatalf("expected water kind, got %q", event.Kind)
			}

			return 101, nil
		},
	}

	svc := NewCareEventService(repo, plantRepo)

	event, err := svc.MarkWater(context.Background(), 77, 15)
	mustNoErr(t, err)

	if event.ID != 101 {
		t.Fatalf("expected event ID 101, got %d", event.ID)
	}
	if event.Kind != domain.CareKindWater {
		t.Fatalf("expected water kind, got %q", event.Kind)
	}
}

func TestCareEventServiceMarkWaterReturnsNotFoundWhenPlantMissing(t *testing.T) {
	createCalled := false
	plantRepo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			return domain.Plant{}, domain.ErrNotFound
		},
	}

	repo := &fakeCareEventRepo{
		createCareEventFn: func(ctx context.Context, event domain.CareEvent) (int64, error) {
			createCalled = true
			return 0, nil
		},
	}

	svc := NewCareEventService(repo, plantRepo)

	_, err := svc.MarkWater(context.Background(), 77, 15)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if createCalled {
		t.Fatal("CreateCareEvent should not be called when plant does not exist")
	}
}

func TestCareEventServiceListCareEventsByTypeChecksPlantOwnership(t *testing.T) {
	plantRepo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			return domain.Plant{ID: plantID, UserID: userID, Name: "Monstera"}, nil
		},
	}

	now := time.Now().UTC()
	repo := &fakeCareEventRepo{
		listCareEventsByTypeFn: func(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error) {
			if plantID != 15 {
				t.Fatalf("expected plant ID 15, got %d", plantID)
			}
			if eventType != domain.CareKindFertilize {
				t.Fatalf("expected fertilize kind, got %q", eventType)
			}

			return []domain.CareEvent{
				{ID: 1, PlantID: 15, Kind: domain.CareKindFertilize, OccurredAt: now},
			}, nil
		},
	}

	svc := NewCareEventService(repo, plantRepo)

	events, err := svc.ListCareEventsByType(context.Background(), 77, 15, domain.CareKindFertilize)
	mustNoErr(t, err)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != domain.CareKindFertilize {
		t.Fatalf("expected fertilize kind, got %q", events[0].Kind)
	}
}

func TestCareEventServiceListCareEventsByTypeReturnsCopy(t *testing.T) {
	plantRepo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			return domain.Plant{ID: plantID, UserID: userID, Name: "Monstera"}, nil
		},
	}

	now := time.Now().UTC()
	repo := &fakeCareEventRepo{
		listCareEventsByTypeFn: func(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error) {
			return []domain.CareEvent{
				{ID: 1, PlantID: 15, Kind: domain.CareKindWater, OccurredAt: now},
			}, nil
		},
	}

	svc := NewCareEventService(repo, plantRepo)

	events, err := svc.ListCareEventsByType(context.Background(), 77, 15, domain.CareKindWater)
	mustNoErr(t, err)

	events[0].Kind = domain.CareKindFertilize

	freshEvents, err := svc.ListCareEventsByType(context.Background(), 77, 15, domain.CareKindWater)
	mustNoErr(t, err)

	if freshEvents[0].Kind != domain.CareKindWater {
		t.Fatalf("expected original kind %q, got %q", domain.CareKindWater, freshEvents[0].Kind)
	}
}

func TestCareEventServiceListRecentCareEventsByTypeSortsAndLimits(t *testing.T) {
	repo := &fakeCareEventRepo{
		listRecentCareEventsByUserAndType: func(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
			if userID != 77 {
				t.Fatalf("expected user ID 77, got %d", userID)
			}
			if eventType != domain.CareKindWater {
				t.Fatalf("expected water kind, got %q", eventType)
			}
			if limit != 3 {
				t.Fatalf("expected limit 3, got %d", limit)
			}

			return []domain.CareEvent{
				{ID: 3, PlantID: 2, Kind: domain.CareKindWater, OccurredAt: time.Date(2026, 4, 17, 11, 0, 0, 0, time.UTC)},
				{ID: 1, PlantID: 1, Kind: domain.CareKindWater, OccurredAt: time.Date(2026, 4, 17, 10, 0, 0, 0, time.UTC)},
				{ID: 4, PlantID: 2, Kind: domain.CareKindWater, OccurredAt: time.Date(2026, 4, 17, 9, 0, 0, 0, time.UTC)},
			}, nil
		},
	}

	svc := NewCareEventService(repo, &fakePlantRepo{})

	events, err := svc.ListRecentCareEventsByType(context.Background(), 77, domain.CareKindWater, 3)
	mustNoErr(t, err)

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	if events[0].ID != 3 || events[1].ID != 1 || events[2].ID != 4 {
		t.Fatalf("unexpected order: %+v", events)
	}
}
