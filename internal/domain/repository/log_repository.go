package repository

import (
	"context"

	"findMyPhone/internal/domain"
)

// LogRepository defines persistence for logs.
type LogRepository interface {
	Create(ctx context.Context, log *domain.Log) error
	GetLastByDeviceID(ctx context.Context, deviceID string) (*domain.Log, error)
}
