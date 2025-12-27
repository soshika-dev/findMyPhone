package repository

import (
	"context"

	"findMyPhone/internal/domain"
)

// DeviceRepository defines persistence for devices.
type DeviceRepository interface {
	Create(ctx context.Context, device *domain.Device) error
	GetByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error)
}
