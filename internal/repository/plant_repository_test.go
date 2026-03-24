package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func newPlantRepositoryTestDB(t *testing.T) (*PlantRepository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	t.Cleanup(func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("db.Close: %v", err)
		}
	})

	return NewPlantRepository(db), mock
}

func TestPlantRepositoryCreatePlantReturnsID(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	createdAt := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)
	plant := domain.Plant{
		UserID:    10,
		Name:      "Monstera",
		CreatedAt: createdAt,
	}

	query := regexp.QuoteMeta(`
		insert into plants (user_id, name, created_at)
		values ($1, $2, $3)
		returning id
	`)

	mock.ExpectQuery(query).
		WithArgs(plant.UserID, plant.Name, plant.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(42)))

	id, err := repo.CreatePlant(ctx, plant)
	if err != nil {
		t.Fatalf("CreatePlant returned error: %v", err)
	}

	if id != 42 {
		t.Fatalf("expected id 42, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPlantRepositoryListPlantsByUserReadsOnlyUserPlants(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(10)
	firstCreatedAt := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)
	secondCreatedAt := firstCreatedAt.Add(time.Hour)

	query := regexp.QuoteMeta(`
		select id, user_id, name, created_at
		from plants
		where user_id = $1
		order by id
	`)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
		AddRow(int64(1), userID, "Monstera", firstCreatedAt).
		AddRow(int64(2), userID, "Cactus", secondCreatedAt)

	mock.ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(rows)

	plants, err := repo.ListPlantsByUser(ctx, userID)
	if err != nil {
		t.Fatalf("ListPlantsByUser returned error: %v", err)
	}

	if len(plants) != 2 {
		t.Fatalf("expected 2 plants, got %d", len(plants))
	}

	for _, plant := range plants {
		if plant.UserID != userID {
			t.Fatalf("expected only user %d plants, got user %d", userID, plant.UserID)
		}
	}

	if plants[0].Name != "Monstera" || plants[1].Name != "Cactus" {
		t.Fatalf("unexpected plants returned: %+v", plants)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPlantRepositoryDeletePlantReturnsNotFoundWhenNoRowsDeleted(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(10)
	plantID := int64(77)

	query := regexp.QuoteMeta(`
		delete from plants
		where id = $1 and user_id = $2
	`)

	mock.ExpectExec(query).
		WithArgs(plantID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeletePlant(ctx, userID, plantID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
