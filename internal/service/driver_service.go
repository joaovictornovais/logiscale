package service

import (
	"context"
	"time"

	"github.com/joaovictornovais/logiscale/internal/domain"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *domain.Driver) error
	GetDriverByID(ctx context.Context, id string) (*domain.Driver, error)
}

type DriverService struct {
	repo DriverRepository
}

func NewDriverService(repo DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}

func (s *DriverService) CreateDriver(ctx context.Context, name string, license string) (*domain.Driver, error) {
	if name == "" || license == "" {
		return nil, domain.ErrInvalidInput
	}

	driver := &domain.Driver{
		Name:      name,
		License:   license,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateDriver(ctx, driver); err != nil {
		return nil, err
	}

	return driver, nil
}

func (s *DriverService) GetDriverByID(ctx context.Context, id string) (*domain.Driver, error) {
	driver, err := s.repo.GetDriverByID(ctx, id)
	if err != nil {
		return nil, domain.ErrDriverNotFound
	}

	return driver, nil
}
