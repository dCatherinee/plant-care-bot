package usecase

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
	"github.com/dCatherinee/plant-care-bot/internal/repo"
)

type PlantService struct {
	repo repo.PlantRepository
}

func NewPlantService(r repo.PlantRepository) *PlantService {
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

	plantID, err := s.repo.CreatePlant(ctx, plant)

	if err != nil {
		return domain.Plant{}, err
	}

	plant.ID = plantID

	return plant, nil
}

func (s *PlantService) ListPlants(ctx context.Context, userID int64) ([]domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	userPlants, err := s.repo.ListPlantsByUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Plant, len(userPlants))
	copy(result, userPlants)
	return result, nil
}

func (s *PlantService) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := s.repo.DeletePlant(ctx, userID, plantID)

	if err != nil {
		return err
	}

	return nil
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
