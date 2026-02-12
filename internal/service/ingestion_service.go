package service

import (
	"context"
	"log"
	"sync"

	"github.com/joaovictornovais/logiscale/internal/domain"
	"github.com/redis/go-redis/v9"
)

type LocationRepository interface {
	SaveLocation(ctx context.Context, loc domain.LocationPayload) error
}

type IngestionService struct {
	repo          LocationRepository
	redis         *redis.Client
	locationQueue chan domain.LocationPayload
	wg            sync.WaitGroup
}

const (
	WorkerCount = 10
	QueueSize   = 1000
)

func NewIngestionService(repo LocationRepository, redisClient *redis.Client) *IngestionService {
	queue := make(chan domain.LocationPayload, QueueSize)

	s := &IngestionService{
		repo:          repo,
		redis:         redisClient,
		locationQueue: queue,
	}

	s.StartWorkers()

	return s
}

func (s *IngestionService) StartWorkers() {
	for i := 0; i < WorkerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *IngestionService) worker(id int) {
	defer s.wg.Done()
	log.Printf("Worker %d started", id)

	for loc := range s.locationQueue {
		ctx := context.Background()

		key := "drivers:locations"
		err := s.redis.GeoAdd(ctx, key, &redis.GeoLocation{
			Name:      loc.DriverID,
			Longitude: loc.Lng,
			Latitude:  loc.Lat,
		}).Err()

		if err != nil {
			log.Printf("[Worker %d] error on regis geoadd: %v", id, err)
		}

		if err := s.repo.SaveLocation(ctx, loc); err != nil {
			log.Printf("[Worker %d] error while trying to save driver's location on postgres %s: %v", id, loc.DriverID, err)
		}
	}
	log.Printf("Worker %d has stopped", id)
}

func (s *IngestionService) Ingest(loc domain.LocationPayload) {
	s.locationQueue <- loc
}

func (s *IngestionService) Close() {
	close(s.locationQueue)
	s.wg.Wait()
}
