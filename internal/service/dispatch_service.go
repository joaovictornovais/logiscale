package service

import (
	"context"
	"fmt"

	"github.com/joaovictornovais/logiscale/internal/domain"
	"github.com/redis/go-redis/v9"
)

type DispatchService struct {
	redis *redis.Client
}

func NewDispatchService(redis *redis.Client) *DispatchService {
	return &DispatchService{redis: redis}
}

func (s *DispatchService) FindNearestDrivers(ctx context.Context, lat, lng, radiusKm float64) ([]domain.RouteResult, error) {
	key := "drivers:locations"

	res, err := s.redis.GeoSearchLocation(ctx, key, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     radiusKm,
			RadiusUnit: "km",
			Sort:       "ASC",
			Count:      10,
		},
		WithDist:  true,
		WithCoord: true,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("error while looking for drivers: %w", err)
	}

	var drivers []domain.RouteResult

	for _, loc := range res {
		drivers = append(drivers, domain.RouteResult{
			DriverID:      loc.Name,
			TotalDistance: loc.Dist,
		})
	}

	return drivers, nil
}
