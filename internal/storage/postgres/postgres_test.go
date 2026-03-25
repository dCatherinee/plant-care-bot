//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/dCatherinee/plant-care-bot/internal/config"
)

var testDB *sql.DB

const testTimeout = 5 * time.Second

func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load test config: %v\n", err)
		os.Exit(1)
	}

	db, err := open(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open test db: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	if err := db.PingContext(ctx); err != nil {
		cancel()
		_ = db.Close()
		fmt.Fprintf(os.Stderr, "ping test db: %v\n", err)
		os.Exit(1)
	}
	cancel()

	if err := applyMigrationsWithDB(db); err != nil {
		_ = db.Close()
		fmt.Fprintf(os.Stderr, "apply migrations: %v\n", err)
		os.Exit(1)
	}

	testDB = db

	code := m.Run()

	if err := db.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close test db: %v\n", err)
		os.Exit(1)
	}

	os.Exit(code)
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	if testDB == nil {
		t.Fatal("test database is not initialized")
	}

	return testDB
}

func applyMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	if err := applyMigrationsWithDB(db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
}

func applyMigrationsWithDB(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	migrationsDir := filepath.Join("..", "..", "..", "migrations")
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func cleanupTables(t *testing.T, db *sql.DB) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	const query = `
		TRUNCATE TABLE reminders, care_events, plants, users
		RESTART IDENTITY CASCADE
	`

	if _, err := db.ExecContext(ctx, query); err != nil {
		t.Fatalf("cleanup tables: %v", err)
	}
}
