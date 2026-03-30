package app

import (
	"database/sql"

	"github.com/dCatherinee/plant-care-bot/internal/storage/postgres"
	"github.com/dCatherinee/plant-care-bot/internal/usecase"
)

const Version = "0.0.1"

type App struct {
	PlantService *usecase.PlantService
}

func New(db *sql.DB) *App {
	plantRepo := postgres.NewPlantRepository(db)
	plantService := usecase.NewPlantService(plantRepo)

	return &App{
		PlantService: plantService,
	}
}
