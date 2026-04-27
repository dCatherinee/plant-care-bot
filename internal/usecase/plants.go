package usecase

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantStore interface {
	PlantFinder

	CreatePlant(ctx context.Context, plant domain.Plant) (domain.Plant, error)
	DeletePlant(ctx context.Context, userID int64, plantID int64) error
	UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error)
}

type PlantService struct {
	repo PlantStore
}

func NewPlantService(r PlantStore) *PlantService {
	return &PlantService{repo: r}
}

func (s *PlantService) AddPlant(ctx context.Context, userID int64, name string) (domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plant{}, err
	}

	plant, err := domain.NewPlant(userID, name)

	if err != nil {
		return domain.Plant{}, err
	}

	return s.repo.CreatePlant(ctx, plant)
}

func (s *PlantService) ListPlants(ctx context.Context, userID int64) ([]domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.repo.ListPlantsByUser(ctx, userID)
}

func (s *PlantService) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return s.repo.DeletePlant(ctx, userID, plantID)
}

func (s *PlantService) GetPlant(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plant{}, err
	}

	plant, err := s.repo.GetPlantByID(ctx, userID, plantID)

	if err != nil {
		return domain.Plant{}, err
	}

	return plant, nil
}

func (s *PlantService) UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plant{}, err
	}

	normalizedName, err := domain.NormalizePlantName(name)
	if err != nil {
		return domain.Plant{}, err
	}

	plant, err := s.repo.UpdatePlantName(ctx, userID, plantID, normalizedName)
	if err != nil {
		return domain.Plant{}, err
	}

	return plant, nil
}
