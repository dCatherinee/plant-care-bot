package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepositoryEnsureUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()
	createdAt := time.Date(2026, time.April, 15, 10, 0, 0, 0, time.UTC)

	query := regexp.QuoteMeta(`
		insert into users (telegram_user_id)
		values ($1)
		on conflict (telegram_user_id)
		do update set telegram_user_id = excluded.telegram_user_id
		returning id, telegram_user_id, created_at
	`)

	rows := sqlmock.NewRows([]string{"id", "telegram_user_id", "created_at"}).
		AddRow(int64(7), int64(12345), createdAt)

	mock.ExpectQuery(query).
		WithArgs(int64(12345)).
		WillReturnRows(rows)

	user, err := repo.EnsureUser(ctx, 12345)
	if err != nil {
		t.Fatalf("EnsureUser returned error: %v", err)
	}

	if user.ID != 7 {
		t.Fatalf("expected user ID %d, got %d", 7, user.ID)
	}

	if user.TelegramUserID != 12345 {
		t.Fatalf("expected telegram user ID %d, got %d", 12345, user.TelegramUserID)
	}

	if user.CreatedAt != createdAt {
		t.Fatalf("expected createdAt %v, got %v", createdAt, user.CreatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestUserRepositoryEnsureUserReturnsWrappedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	query := regexp.QuoteMeta(`
		insert into users (telegram_user_id)
		values ($1)
		on conflict (telegram_user_id)
		do update set telegram_user_id = excluded.telegram_user_id
		returning id, telegram_user_id, created_at
	`)

	mock.ExpectQuery(query).
		WithArgs(int64(12345)).
		WillReturnError(errors.New("db failed"))

	_, err = repo.EnsureUser(ctx, 12345)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "ensure user: db failed" {
		t.Fatalf("expected wrapped error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
