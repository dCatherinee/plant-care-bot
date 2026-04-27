package postgres

import (
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type plant struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	name       string    `db:"name"`
	created_at time.Time `db:"created_at"`
}

func (p *plant) scanPlant() []any {
	return []any{
		&p.ID,
		&p.UserID,
		&p.name,
		&p.created_at,
	}
}

func newPlant(value plant) domain.Plant {
	// Plants are validated before persistence. If a row can no longer be
	// reconstructed into a valid domain entity, we treat it as a broken storage
	// invariant and fail fast instead of masking the corruption.
	return domain.MustPlant(int64(value.ID), int64(value.UserID), value.name, value.created_at)
}

func newPlants(values []plant) []domain.Plant {
	res := make([]domain.Plant, len(values))

	for i, value := range values {
		res[i] = newPlant(value)
	}

	return res
}
