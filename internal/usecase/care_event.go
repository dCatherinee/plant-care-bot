package usecase

import (
	"context"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type CareEventStore interface {
	CreateCareEvent(ctx context.Context, event domain.CareEvent) (domain.CareEvent, error)
	ListCareEventsByType(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error)
	ListRecentCareEventsByUserAndType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error)
}

type CareEventService struct {
	repo      CareEventStore
	plantRepo PlantFinder
}

func NewCareEventService(r CareEventStore, plantRepo PlantFinder) *CareEventService {
	return &CareEventService{repo: r, plantRepo: plantRepo}
}

func (s *CareEventService) ensurePlantExists(ctx context.Context, userID int64, plantID int64) error {
	if _, err := s.plantRepo.GetPlantByID(ctx, userID, plantID); err != nil {
		return err
	}

	return nil
}

func (s *CareEventService) markCare(ctx context.Context, userID int64, plantID int64, careType domain.CareKind) (domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return domain.CareEvent{}, err
	}

	if err := s.ensurePlantExists(ctx, userID, plantID); err != nil {
		return domain.CareEvent{}, err
	}

	careEvent, err := domain.NewCareEvent(plantID, careType, time.Now())

	if err != nil {
		return domain.CareEvent{}, err
	}

	return s.repo.CreateCareEvent(ctx, careEvent)
}

func (s *CareEventService) MarkWater(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return domain.CareEvent{}, err
	}

	return s.markCare(ctx, userID, plantID, domain.CareKindWater)
}

func (s *CareEventService) MarkFertilize(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return domain.CareEvent{}, err
	}

	return s.markCare(ctx, userID, plantID, domain.CareKindFertilize)
}

func (s *CareEventService) ListCareEventsByType(ctx context.Context, userID int64, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := s.ensurePlantExists(ctx, userID, plantID); err != nil {
		return nil, err
	}

	return s.repo.ListCareEventsByType(ctx, plantID, eventType)
}

func (s *CareEventService) ListRecentCareEventsByType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.repo.ListRecentCareEventsByUserAndType(ctx, userID, eventType, limit)
}
