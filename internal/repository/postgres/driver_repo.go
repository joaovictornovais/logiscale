package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictornovais/logiscale/internal/domain"
)

type DriverRepository struct {
	db *pgxpool.Pool
}

func NewDriverRepository(db *pgxpool.Pool) *DriverRepository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) CreateDriver(ctx context.Context, driver *domain.Driver) error {
	query := `
		INSERT INTO drivers (name, license, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		driver.Name,
		driver.License,
		driver.CreatedAt,
	).Scan(&driver.ID)

	if err != nil {
		return fmt.Errorf("error while inserting driver: %w", err)
	}

	return nil
}

func (r *DriverRepository) GetDriverByID(ctx context.Context, id string) (*domain.Driver, error) {
	query := `SELECT id, name, license, created_at FROM drivers WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	var d domain.Driver

	if err := row.Scan(&d.ID, &d.Name, &d.License, &d.CreatedAt); err != nil {
		return nil, fmt.Errorf("error while scanning driver data: %w", err)
	}

	return &d, nil
}
