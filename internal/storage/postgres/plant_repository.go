package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type PlantRepository struct {
	db *sql.DB
}

func NewPlantRepository(db *sql.DB) *PlantRepository {
	return &PlantRepository{db: db}
}

func (r *PlantRepository) CreatePlant(ctx context.Context, value domain.Plant) (domain.Plant, error) {
	var (
		query = `
		insert into plants (user_id, name, created_at)
		values ($1, $2, $3)
		returning id, user_id, name, created_at
	`
		dest plant
	)

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	if err := r.db.QueryRowContext(ctx, query, value.UserID, value.Name, value.CreatedAt).Scan(dest.scanPlant()...); err != nil {
		return domain.Plant{}, fmt.Errorf("create plant: %w", err)
	}

	return newPlant(dest), nil
}

func (r *PlantRepository) ListPlantsByUser(ctx context.Context, userID int64) ([]domain.Plant, error) {
	const query = `
		select id, user_id, name, created_at
		from plants
		where user_id = $1
		order by id
	`

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list plants by user: %w", err)
	}
	defer rows.Close()

	var res []plant
	for rows.Next() {
		var plant plant
		if err := rows.Scan(plant.scanPlant()...); err != nil {
			return nil, fmt.Errorf("scan plant: %w", err)
		}

		res = append(res, plant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plants: %w", err)
	}

	return newPlants(res), nil
}

func (r *PlantRepository) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	const query = `
		delete from plants
		where id = $1 and user_id = $2
	`

	ctx, cancel := withTimeout(ctx)
	defer cancel()

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

func (r *PlantRepository) GetPlantByID(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
	var (
		query = `
			select id, user_id, name, created_at from plants
			where id = $1 and user_id = $2
		`
		dest plant
	)

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, plantID, userID).Scan(dest.scanPlant()...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Plant{}, domain.ErrNotFound
		}
		return domain.Plant{}, fmt.Errorf("get plant by id: %w", err)
	}

	return newPlant(dest), nil
}

func (r *PlantRepository) UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
	const queryUpdate = `
		update plants set name = $1
		where id = $2 and user_id = $3
	`

	const querySelect = `
		select id, user_id, name, created_at from plants
		where id = $1 and user_id = $2
	`

	ctx, cancel := withTimeout(ctx)
	defer cancel()

	result, err := r.db.ExecContext(ctx, queryUpdate, name, plantID, userID)
	if err != nil {
		return domain.Plant{}, fmt.Errorf("update plant name: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.Plant{}, fmt.Errorf("update plant name rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.Plant{}, fmt.Errorf("update plant name: %w", domain.ErrNotFound)
	}

	var dest plant
	err = r.db.QueryRowContext(ctx, querySelect, plantID, userID).Scan(dest.scanPlant()...)

	if err != nil {
		return domain.Plant{}, fmt.Errorf("scan plant: %w", err)
	}

	return newPlant(dest), nil
}
