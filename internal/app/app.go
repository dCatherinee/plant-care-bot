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

type Storage struct {
	Plant     *postgres.PlantRepository
	User      *postgres.UserRepository
	CareEvent *postgres.CareEventRepository
}

type UseCases struct {
	Plants    *usecase.PlantService
	User      *usecase.UserService
	CareEvent *usecase.CareEventService
}

type Container struct {
	UseCases UseCases
	Storage  Storage
}

func New(db *sql.DB) *App {
	storages := Storage{
		Plant:     postgres.NewPlantRepository(db),
		User:      postgres.NewUserRepository(db),
		CareEvent: postgres.NewCareEventRepository(db),
	}

	usecases := UseCases{
		Plants:    usecase.NewPlantService(storages.Plant),
		User:      usecase.NewUserService(storages.User),
		CareEvent: usecase.NewCareEventService(storages.CareEvent, storages.Plant),
	}

	container := Container{
		Storage:  storages,
		UseCases: usecases,
	}

	return &App{
		PlantService:     container.UseCases.Plants,
		UserService:      container.UseCases.User,
		CareEventService: container.UseCases.CareEvent,
	}
}
