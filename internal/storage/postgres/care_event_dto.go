package postgres

import (
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type careEvent struct {
	ID         int             `db:"id"`
	PlantID    int             `db:"plant_id"`
	Kind       domain.CareKind `db:"event_type"`
	OccurredAt time.Time       `db:"occurred_at"`
	CreatedAt  time.Time       `db:"created_at"`
}

func (c *careEvent) scanCareEvent() []any {
	return []any{
		&c.ID,
		&c.PlantID,
		&c.Kind,
		&c.OccurredAt,
		&c.CreatedAt,
	}
}

func newCareEvent(value careEvent) domain.CareEvent {
	// Care events are validated before persistence. If a row can no longer be
	// reconstructed into a valid domain entity, we treat it as a broken storage
	// invariant and fail fast instead of masking the corruption.
	return domain.MustCareEvent(int64(value.ID), int64(value.PlantID), value.Kind, value.OccurredAt, value.CreatedAt)
}

func newCareEvents(values []careEvent) []domain.CareEvent {
	res := make([]domain.CareEvent, len(values))

	for i, value := range values {
		res[i] = newCareEvent(value)
	}

	return res
}
