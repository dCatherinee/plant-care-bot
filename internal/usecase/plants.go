package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantService struct {
	mu           sync.Mutex
	nextID       int64
	plantsByUser map[int64][]domain.Plant
}

func NewPlantService() *PlantService {
	return &PlantService{
		nextID:       0,
		plantsByUser: make(map[int64][]domain.Plant),
	}
}

func (s *PlantService) AddPlant(ctx context.Context, userID int64, name string) (domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plant{}, err
	}

	plant, err := domain.NewPlant(userID, name)

	if err != nil {
		return domain.Plant{}, err
	}

	select {
	case <-time.After(50 * time.Millisecond):
		s.mu.Lock()
		defer s.mu.Unlock()

		s.nextID++
		plant.ID = s.nextID
		s.plantsByUser[userID] = append(s.plantsByUser[userID], plant)

		return plant, nil
	case <-ctx.Done():
		return domain.Plant{}, ctx.Err()
	}
}

func (s *PlantService) ListPlants(ctx context.Context, userID int64) ([]domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	userPlants := s.plantsByUser[userID]

	result := make([]domain.Plant, len(userPlants))
	copy(result, userPlants)
	return result, nil
}

func (s *PlantService) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	userPlants := s.plantsByUser[userID]

	for i, plant := range userPlants {
		if plant.ID != plantID {
			continue
		}

		userPlants = append(userPlants[:i], userPlants[i+1:]...)
		s.plantsByUser[userID] = userPlants
		return nil
	}

	return domain.ErrNotFound
}

func (s *PlantService) GetPlant(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plant{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	userPlants := s.plantsByUser[userID]

	for _, plant := range userPlants {
		if plantID == plant.ID {
			return plant, nil
		}
	}

	return domain.Plant{}, domain.ErrNotFound
}

func (s *PlantService) UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plant{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	userPlants := s.plantsByUser[userID]

	for i := range userPlants {
		if plantID != userPlants[i].ID {
			continue
		}

		if err := userPlants[i].Rename(name); err != nil {
			return domain.Plant{}, err
		}

		return userPlants[i], nil
	}

	return domain.Plant{}, domain.ErrNotFound
}
