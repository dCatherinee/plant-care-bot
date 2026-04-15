package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) EnsureUser(ctx context.Context, telegramUserID int64) (domain.User, error) {
	const query = `
		insert into users (telegram_user_id)
		values ($1)
		on conflict (telegram_user_id)
		do update set telegram_user_id = excluded.telegram_user_id
		returning id, telegram_user_id, created_at
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, telegramUserID).Scan(
		&user.ID,
		&user.TelegramUserID,
		&user.CreatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("ensure user: %w", err)
	}

	return user, nil
}
