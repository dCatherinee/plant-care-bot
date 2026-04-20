package usecase

import (
	"context"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
	"github.com/dCatherinee/plant-care-bot/internal/repo"
)

type CareEventService struct {
	repo      repo.CareEventRepository
	plantRepo repo.PlantRepository
}

func NewCareEventService(r repo.CareEventRepository, plantRepo repo.PlantRepository) *CareEventService {
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

	careEventID, err := s.repo.CreateCareEvent(ctx, careEvent)

	if err != nil {
		return domain.CareEvent{}, err
	}

	careEvent.ID = careEventID

	return careEvent, nil
}

func (s *CareEventService) MarkWater(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return domain.CareEvent{}, err
	}

	careEvent, err := s.markCare(ctx, userID, plantID, domain.CareKindWater)
	if err != nil {
		return domain.CareEvent{}, err
	}

	return careEvent, nil
}

func (s *CareEventService) MarkFertilize(ctx context.Context, userID int64, plantID int64) (domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return domain.CareEvent{}, err
	}

	careEvent, err := s.markCare(ctx, userID, plantID, domain.CareKindFertilize)
	if err != nil {
		return domain.CareEvent{}, err
	}

	return careEvent, nil
}

func (s *CareEventService) ListCareEventsByType(ctx context.Context, userID int64, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := s.ensurePlantExists(ctx, userID, plantID); err != nil {
		return nil, err
	}

	careEvents, err := s.repo.ListCareEventsByType(ctx, plantID, eventType)

	if err != nil {
		return nil, err
	}

	result := make([]domain.CareEvent, len(careEvents))
	copy(result, careEvents)
	return result, nil
}

func (s *CareEventService) ListRecentCareEventsByType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	careEvents, err := s.repo.ListRecentCareEventsByUserAndType(ctx, userID, eventType, limit)
	if err != nil {
		return nil, err
	}

	result := make([]domain.CareEvent, len(careEvents))
	copy(result, careEvents)
	return result, nil
}
