package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictornovais/logiscale/internal/domain"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) SaveLocation(ctx context.Context, loc domain.LocationPayload) error {
	query := `INSERT INTO locations (driver_id, lat, lng, sent_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, loc.DriverID, loc.Lat, loc.Lng, loc.SentAt)
	if err != nil {
		return fmt.Errorf("failed to save location: %w", err)
	}
	return nil
}
