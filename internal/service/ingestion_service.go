package service

import (
	"context"
	"log"
	"sync"

	"github.com/joaovictornovais/logiscale/internal/domain"
)

type LocationRepository interface {
	SaveLocation(ctx context.Context, loc domain.LocationPayload) error
}

type IngestionService struct {
	repo          LocationRepository
	locationQueue chan domain.LocationPayload
	wg            sync.WaitGroup
}

const (
	WorkerCount = 10
	QueueSize   = 1000
)

func NewIngestionService(repo LocationRepository) *IngestionService {
	queue := make(chan domain.LocationPayload, QueueSize)

	s := &IngestionService{
		repo:          repo,
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

		err := s.repo.SaveLocation(ctx, loc)
		if err != nil {
			log.Printf("[Worker %d] error while trying to save driver's location %s: %v", id, loc.DriverID, err)
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
