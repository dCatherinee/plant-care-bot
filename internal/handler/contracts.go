package handler

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantUsecase interface {
	AddPlant(ctx context.Context, userID int64, name string) (domain.Plant, error)
	ListPlants(ctx context.Context, userID int64) ([]domain.Plant, error)
	GetPlant(ctx context.Context, userID int64, plantID int64) (domain.Plant, error)
	UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error)
	DeletePlant(ctx context.Context, userID int64, plantID int64) error
}
