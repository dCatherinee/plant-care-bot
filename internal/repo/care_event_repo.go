package repo

import (
	"context"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type CareEventRepository interface {
	CreateCareEvent(ctx context.Context, event domain.CareEvent) (int64, error)
	ListCareEventsByType(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error)
	ListRecentCareEventsByUserAndType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error)
}
