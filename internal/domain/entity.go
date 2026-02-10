package domain

import (
	"errors"
	"time"
)

var (
	ErrDriverNotFound = errors.New("driver not found")
	ErrInvalidCoord   = errors.New("invalid coordinates")
	ErrInvalidInput   = errors.New("invalid driver data")
)

type Driver struct {
	ID        string
	Name      string
	License   string
	CreatedAt time.Time
}

type TruckLocation struct {
	TruckID   string
	DriverID  string
	Lat       float64
	Lng       float64
	Timestamp time.Time
}

type RouteResult struct {
	DriverID      string
	TotalDistance float64
	EstimatedTime int
	PathPoints    []string
}
