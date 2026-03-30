package repo

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantRepository interface {
	CreatePlant(ctx context.Context, plant domain.Plant) (int64, error)
	ListPlantsByUser(ctx context.Context, userID int64) ([]domain.Plant, error)
	DeletePlant(ctx context.Context, userID int64, plantID int64) error
	GetPlantByID(ctx context.Context, userID int64, plantID int64) (domain.Plant, error)
	UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error)
}
