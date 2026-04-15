package repo

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type UserRepository interface {
	EnsureUser(ctx context.Context, telegramUserID int64) (domain.User, error)
}
