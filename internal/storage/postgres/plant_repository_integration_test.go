//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func createTestUser(t *testing.T, telegramUserID int64) int64 {
	t.Helper()

	db := newTestDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	const query = `
		insert into users (telegram_user_id)
		values ($1)
		returning id
	`

	var userID int64
	if err := db.QueryRowContext(ctx, query, telegramUserID).Scan(&userID); err != nil {
		t.Fatalf("create test user: %v", err)
	}

	return userID
}

func createTestPlant(t *testing.T, userID int64, name string, createdAt time.Time) int64 {
	t.Helper()

	repo := NewPlantRepository(newTestDB(t))
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	id, err := repo.CreatePlant(ctx, domain.Plant{
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
	})
	if err != nil {
		t.Fatalf("create test plant: %v", err)
	}

	return id
}

func TestPlantRepositoryCreatePlant_Integration(t *testing.T) {
	db := newTestDB(t)
	cleanupTables(t, db)

	userID := createTestUser(t, 2001)
	repo := NewPlantRepository(db)
	ctx := context.Background()
	createdAt := time.Date(2026, time.March, 25, 12, 0, 0, 0, time.UTC)

	plant := domain.Plant{
		UserID:    userID,
		Name:      "Monstera",
		CreatedAt: createdAt,
	}

	id, err := repo.CreatePlant(ctx, plant)
	if err != nil {
		t.Fatalf("CreatePlant returned error: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected positive id, got %d", id)
	}

	const query = `
		select id, user_id, name, created_at
		from plants
		where id = $1
	`

	var saved domain.Plant
	if err := db.QueryRowContext(ctx, query, id).Scan(
		&saved.ID,
		&saved.UserID,
		&saved.Name,
		&saved.CreatedAt,
	); err != nil {
		t.Fatalf("query saved plant: %v", err)
	}

	if saved.ID != id {
		t.Fatalf("expected saved id %d, got %d", id, saved.ID)
	}
	if saved.UserID != plant.UserID {
		t.Fatalf("expected user id %d, got %d", plant.UserID, saved.UserID)
	}
	if saved.Name != plant.Name {
		t.Fatalf("expected name %q, got %q", plant.Name, saved.Name)
	}
	if !saved.CreatedAt.Equal(plant.CreatedAt) {
		t.Fatalf("expected created_at %v, got %v", plant.CreatedAt, saved.CreatedAt)
	}
}

func TestPlantRepositoryListPlantsByUser_Integration(t *testing.T) {
	db := newTestDB(t)
	cleanupTables(t, db)

	firstUserID := createTestUser(t, 2002)
	secondUserID := createTestUser(t, 2003)

	createdAt := time.Date(2026, time.March, 25, 13, 0, 0, 0, time.UTC)
	firstPlantID := createTestPlant(t, firstUserID, "Monstera", createdAt)
	secondPlantID := createTestPlant(t, firstUserID, "Cactus", createdAt.Add(time.Hour))
	createTestPlant(t, secondUserID, "Orchid", createdAt.Add(2*time.Hour))

	repo := NewPlantRepository(db)
	ctx := context.Background()

	plants, err := repo.ListPlantsByUser(ctx, firstUserID)
	if err != nil {
		t.Fatalf("ListPlantsByUser returned error: %v", err)
	}

	if len(plants) != 2 {
		t.Fatalf("expected 2 plants, got %d", len(plants))
	}

	if plants[0].ID != firstPlantID {
		t.Fatalf("expected first plant id %d, got %d", firstPlantID, plants[0].ID)
	}
	if plants[0].UserID != firstUserID {
		t.Fatalf("expected first plant user id %d, got %d", firstUserID, plants[0].UserID)
	}
	if plants[0].Name != "Monstera" {
		t.Fatalf("expected first plant name %q, got %q", "Monstera", plants[0].Name)
	}

	if plants[1].ID != secondPlantID {
		t.Fatalf("expected second plant id %d, got %d", secondPlantID, plants[1].ID)
	}
	if plants[1].UserID != firstUserID {
		t.Fatalf("expected second plant user id %d, got %d", firstUserID, plants[1].UserID)
	}
	if plants[1].Name != "Cactus" {
		t.Fatalf("expected second plant name %q, got %q", "Cactus", plants[1].Name)
	}
}

func TestPlantRepositoryDeletePlant_Integration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		userID := createTestUser(t, 2004)
		plantID := createTestPlant(t, userID, "Monstera", time.Date(2026, time.March, 25, 15, 0, 0, 0, time.UTC))
		repo := NewPlantRepository(db)
		ctx := context.Background()

		if err := repo.DeletePlant(ctx, userID, plantID); err != nil {
			t.Fatalf("DeletePlant returned error: %v", err)
		}

		const query = `
			select id
			from plants
			where id = $1
		`

		var deletedPlantID int64
		err := db.QueryRowContext(ctx, query, plantID).Scan(&deletedPlantID)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows after delete, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		userID := createTestUser(t, 2005)
		repo := NewPlantRepository(db)
		ctx := context.Background()

		err := repo.DeletePlant(ctx, userID, 9999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("wrong_user", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		ownerUserID := createTestUser(t, 2006)
		otherUserID := createTestUser(t, 2007)
		plantID := createTestPlant(t, ownerUserID, "Monstera", time.Date(2026, time.March, 25, 16, 0, 0, 0, time.UTC))
		repo := NewPlantRepository(db)
		ctx := context.Background()

		err := repo.DeletePlant(ctx, otherUserID, plantID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPlantRepositoryGetPlantByID_Integration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		userID := createTestUser(t, 2004)
		plantID := createTestPlant(t, userID, "Monstera", time.Date(2026, time.March, 25, 15, 0, 0, 0, time.UTC))
		repo := NewPlantRepository(db)
		ctx := context.Background()

		if _, err := repo.GetPlantByID(ctx, userID, plantID); err != nil {
			t.Fatalf("GetPlantByID returned error: %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		userID := createTestUser(t, 2004)
		repo := NewPlantRepository(db)
		ctx := context.Background()

		_, err := repo.GetPlantByID(ctx, userID, 10)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("db_failed", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		db, err := sql.Open("pgx", "postgres://invalid:invalid@127.0.0.1:1/invalid?sslmode=disable")
		if err != nil {
			t.Fatalf("sql.Open returned error: %v", err)
		}
		defer db.Close()

		repo := NewPlantRepository(db)

		_, err = repo.GetPlantByID(ctx, 10, 1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected db error, got ErrNotFound: %v", err)
		}
	})
}

func TestPlantRepositoryUpdatePlantName_Integration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		userID := createTestUser(t, 2004)
		plantID := createTestPlant(t, userID, "Monstera", time.Date(2026, time.March, 25, 15, 0, 0, 0, time.UTC))
		newName := "Cactus"
		repo := NewPlantRepository(db)
		ctx := context.Background()

		plant, err := repo.UpdatePlantName(ctx, userID, plantID, newName)
		if err != nil {
			t.Fatalf("UpdatePlantName returned error: %v", err)
		}

		if plant.Name != newName {
			t.Fatalf("expected plant name %q, got %q", newName, plant.Name)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		userID := createTestUser(t, 2004)
		newName := "Cactus"
		repo := NewPlantRepository(db)
		ctx := context.Background()

		_, err := repo.UpdatePlantName(ctx, userID, 10, newName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("wrong_user", func(t *testing.T) {
		db := newTestDB(t)
		cleanupTables(t, db)

		firstUserID := createTestUser(t, 2004)
		secondUserID := createTestUser(t, 2005)
		plantID := createTestPlant(t, firstUserID, "Monstera", time.Date(2026, time.March, 25, 15, 0, 0, 0, time.UTC))
		newName := "Cactus"
		repo := NewPlantRepository(db)
		ctx := context.Background()

		_, err := repo.UpdatePlantName(ctx, secondUserID, plantID, newName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}
