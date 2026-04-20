package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type CareEventRepository struct {
	db *sql.DB
}

func NewCareEventRepository(db *sql.DB) *CareEventRepository {
	return &CareEventRepository{db: db}
}

func (r *CareEventRepository) CreateCareEvent(ctx context.Context, event domain.CareEvent) (int64, error) {
	const query = `
		insert into care_events (plant_id, event_type, occurred_at)
		values ($1, $2, $3)
		returning id;
	`

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	var id int64
	if err := r.db.QueryRowContext(ctx, query, event.PlantID, event.Kind, event.OccurredAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("create care event: %w", err)
	}

	return id, nil
}

func (r *CareEventRepository) ListCareEventsByType(ctx context.Context, plantID int64, eventType domain.CareKind) ([]domain.CareEvent, error) {
	const query = `
		select id, plant_id, event_type, occurred_at, created_at
		from care_events
		where plant_id = $1 and event_type = $2
		order by occurred_at desc
	`

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, plantID, eventType)
	if err != nil {
		return nil, fmt.Errorf("list care event by user: %w", err)
	}
	defer rows.Close()

	var careEvents []domain.CareEvent
	for rows.Next() {
		var event domain.CareEvent
		if err := rows.Scan(&event.ID, &event.PlantID, &event.Kind, &event.OccurredAt, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan plant: %w", err)
		}

		careEvents = append(careEvents, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plants: %w", err)
	}

	return careEvents, nil
}

func (r *CareEventRepository) ListRecentCareEventsByUserAndType(ctx context.Context, userID int64, eventType domain.CareKind, limit int) ([]domain.CareEvent, error) {
	const query = `
		select ce.id, ce.plant_id, ce.event_type, ce.occurred_at, ce.created_at
		from care_events ce
		join plants p on p.id = ce.plant_id
		where p.user_id = $1 and ce.event_type = $2
		order by ce.occurred_at desc
		limit $3
	`

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID, eventType, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent care events by user and type: %w", err)
	}
	defer rows.Close()

	var careEvents []domain.CareEvent
	for rows.Next() {
		var event domain.CareEvent
		if err := rows.Scan(&event.ID, &event.PlantID, &event.Kind, &event.OccurredAt, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan care event: %w", err)
		}

		careEvents = append(careEvents, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate care events: %w", err)
	}

	return careEvents, nil
}
