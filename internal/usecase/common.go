package usecase

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantFinder interface {
	ListPlantsByUser(ctx context.Context, userID int64) ([]domain.Plant, error)
	GetPlantByID(ctx context.Context, userID int64, plantID int64) (domain.Plant, error)
}
