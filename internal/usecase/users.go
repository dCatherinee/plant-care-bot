package usecase

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
	"github.com/dCatherinee/plant-care-bot/internal/repo"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) EnsureUser(ctx context.Context, telegramUserID int64) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, err
	}

	if _, err := domain.NewUser(telegramUserID); err != nil {
		return domain.User{}, err
	}

	user, err := s.repo.EnsureUser(ctx, telegramUserID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
