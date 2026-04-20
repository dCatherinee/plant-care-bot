package app

import (
	"database/sql"

	"github.com/dCatherinee/plant-care-bot/internal/storage/postgres"
	"github.com/dCatherinee/plant-care-bot/internal/usecase"
)

const Version = "0.0.1"

type App struct {
	PlantService     *usecase.PlantService
	UserService      *usecase.UserService
	CareEventService *usecase.CareEventService
}

func New(db *sql.DB) *App {
	plantRepo := postgres.NewPlantRepository(db)
	userRepo := postgres.NewUserRepository(db)
	careEventRepo := postgres.NewCareEventRepository(db)
	plantService := usecase.NewPlantService(plantRepo)
	userService := usecase.NewUserService(userRepo)
	careEventService := usecase.NewCareEventService(careEventRepo, plantRepo)

	return &App{
		PlantService:     plantService,
		UserService:      userService,
		CareEventService: careEventService,
	}
}
