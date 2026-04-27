package postgres

import (
	"context"
	"database/sql"
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

func TestPlantRepositoryCreatePlantReturnsPlant(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	createdAt := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)
	input := domain.Plant{
		UserID:    10,
		Name:      "Monstera",
		CreatedAt: createdAt,
	}
	expected := domain.Plant{
		ID:        42,
		UserID:    10,
		Name:      "Monstera",
		CreatedAt: createdAt,
	}

	query := regexp.QuoteMeta(`
		insert into plants (user_id, name, created_at)
		values ($1, $2, $3)
		returning id, user_id, name, created_at
	`)

	mock.ExpectQuery(query).
		WithArgs(input.UserID, input.Name, input.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
			AddRow(expected.ID, expected.UserID, expected.Name, expected.CreatedAt))

	actual, err := repo.CreatePlant(ctx, input)
	if err != nil {
		t.Fatalf("CreatePlant returned error: %v", err)
	}

	if actual != expected {
		t.Fatalf("expected plant %+v, got %+v", expected, actual)
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

func TestPlantRepositoryGetPlantByIDReturnCorrectPlant(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	createdAt := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)
	plantID := int64(1)
	expected := domain.Plant{
		UserID:    10,
		Name:      "Monstera",
		CreatedAt: createdAt,
	}

	query := regexp.QuoteMeta(`
		select id, user_id, name, created_at
		from plants
		where id = $1 and user_id = $2
	`)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
		AddRow(plantID, expected.UserID, expected.Name, expected.CreatedAt)

	mock.ExpectQuery(query).
		WithArgs(plantID, expected.UserID).
		WillReturnRows(rows)

	actual, err := repo.GetPlantByID(ctx, expected.UserID, plantID)
	if err != nil {
		t.Fatalf("GetPlantByID returned error: %v", err)
	}

	if actual.ID != plantID {
		t.Fatalf("unexpected plant returned: %v", actual)
	}

	if actual.UserID != expected.UserID {
		t.Fatalf("unexpected plant returned: %v", actual)
	}

	if actual.CreatedAt != expected.CreatedAt {
		t.Fatalf("unexpected plant returned: %v", actual)
	}

	if actual.Name != expected.Name {
		t.Fatalf("unexpected plant returned: %v", actual)
	}
}

func TestPlantRepositoryGetPlantByIDReturnsNotFound(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(10)
	plantID := int64(77)

	query := regexp.QuoteMeta(`
		select id, user_id, name, created_at
		from plants
		where id = $1 and user_id = $2
	`)

	mock.ExpectQuery(query).
		WithArgs(plantID, userID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetPlantByID(ctx, userID, plantID)
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

func TestPlantRepositoryGetPlantByIDReturnsDbFailed(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(10)
	plantID := int64(77)

	query := regexp.QuoteMeta(`
		select id, user_id, name, created_at
		from plants
		where id = $1 and user_id = $2
	`)

	mock.ExpectQuery(query).
		WithArgs(plantID, userID).
		WillReturnError(errors.New("db failed"))

	_, err := repo.GetPlantByID(ctx, userID, plantID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected db failed, got ErrNotFound: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPlantRepositoryUpdatePlantNameChangeName(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	createdAt := time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC)
	userID := int64(10)
	plantID := int64(1)
	newName := "Cactus"

	queryUpdate := regexp.QuoteMeta(`
		update plants set name = $1
		where id = $2 and user_id = $3
	`)

	querySelect := regexp.QuoteMeta(`
		select id, user_id, name, created_at from plants
		where id = $1 and user_id = $2
	`)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
		AddRow(plantID, userID, newName, createdAt)

	mock.ExpectExec(queryUpdate).
		WithArgs(newName, plantID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(querySelect).
		WithArgs(plantID, userID).
		WillReturnRows(rows)

	actual, err := repo.UpdatePlantName(ctx, userID, plantID, newName)
	if err != nil {
		t.Fatalf("UpdatePlantName returned error: %v", err)
	}

	if actual.ID != plantID {
		t.Fatalf("expected plant ID %d, got %d", plantID, actual.ID)
	}

	if actual.UserID != userID {
		t.Fatalf("expected user ID %d, got %d", userID, actual.UserID)
	}

	if actual.Name != newName {
		t.Fatalf("expected plant name %q, got %q", newName, actual.Name)
	}

	if actual.CreatedAt != createdAt {
		t.Fatalf("expected createdAt %v, got %v", createdAt, actual.CreatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestPlantRepositoryUpdatePlantNameReturnsNotFoundWhenNoRowsUpdated(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(10)
	plantID := int64(77)

	queryUpdate := regexp.QuoteMeta(`
		update plants set name = $1
		where id = $2 and user_id = $3
	`)

	mock.ExpectExec(queryUpdate).
		WithArgs("Cactus", plantID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := repo.UpdatePlantName(ctx, userID, plantID, "Cactus")
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

func TestPlantRepositoryUpdatePlantNameReturnsSelectError(t *testing.T) {
	repo, mock := newPlantRepositoryTestDB(t)
	ctx := context.Background()
	userID := int64(10)
	plantID := int64(1)
	newName := "Cactus"

	queryUpdate := regexp.QuoteMeta(`
		update plants set name = $1
		where id = $2 and user_id = $3
	`)

	querySelect := regexp.QuoteMeta(`
		select id, user_id, name, created_at from plants
		where id = $1 and user_id = $2
	`)

	mock.ExpectExec(queryUpdate).
		WithArgs(newName, plantID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(querySelect).
		WithArgs(plantID, userID).
		WillReturnError(errors.New("select failed"))

	_, err := repo.UpdatePlantName(ctx, userID, plantID, newName)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "scan plant: select failed" {
		t.Fatalf("expected wrapped select error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
