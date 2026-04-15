package app

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestVersion_NotEmpty(t *testing.T) {
	if Version == "" {
		t.Fatal("Version must not be empty")
	}
}

func TestNewBuildsServices(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	app := New(db)

	if app.PlantService == nil {
		t.Fatal("expected PlantService to be initialized")
	}

	if app.UserService == nil {
		t.Fatal("expected UserService to be initialized")
	}
}
