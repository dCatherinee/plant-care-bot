package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantRepository struct {
	db *sql.DB
}

func NewPlantRepository(db *sql.DB) *PlantRepository {
	return &PlantRepository{db: db}
}

func (r *PlantRepository) CreatePlant(ctx context.Context, plant domain.Plant) (int64, error) {
	const query = `
		insert into plants (user_id, name, created_at)
		values ($1, $2, $3)
		returning id
	`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, plant.UserID, plant.Name, plant.CreatedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("create plant: %w", err)
	}

	return id, nil
}

func (r *PlantRepository) ListPlantsByUser(ctx context.Context, userID int64) ([]domain.Plant, error) {
	const query = `
		select id, user_id, name, created_at
		from plants
		where user_id = $1
		order by id
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list plants by user: %w", err)
	}
	defer rows.Close()

	var plants []domain.Plant
	for rows.Next() {
		var plant domain.Plant
		if err := rows.Scan(&plant.ID, &plant.UserID, &plant.Name, &plant.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan plant: %w", err)
		}

		plants = append(plants, plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plants: %w", err)
	}

	return plants, nil
}

func (r *PlantRepository) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	const query = `
		delete from plants
		where id = $1 and user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, plantID, userID)
	if err != nil {
		return fmt.Errorf("delete plant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete plant rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("delete plant: %w", domain.ErrNotFound)
	}

	return nil
}
