//go:build integration

package postgres

import (
	"context"
	"testing"
)

func TestUserRepositoryEnsureUser_Integration(t *testing.T) {
	db := newTestDB(t)
	cleanupTables(t, db)

	repo := NewUserRepository(db)
	ctx := context.Background()

	user, err := repo.EnsureUser(ctx, 3001)
	if err != nil {
		t.Fatalf("EnsureUser returned error: %v", err)
	}

	if user.ID <= 0 {
		t.Fatalf("expected positive user ID, got %d", user.ID)
	}

	if user.TelegramUserID != 3001 {
		t.Fatalf("expected telegram user ID %d, got %d", 3001, user.TelegramUserID)
	}

	sameUser, err := repo.EnsureUser(ctx, 3001)
	if err != nil {
		t.Fatalf("EnsureUser second call returned error: %v", err)
	}

	if sameUser.ID != user.ID {
		t.Fatalf("expected same user ID %d, got %d", user.ID, sameUser.ID)
	}
}
