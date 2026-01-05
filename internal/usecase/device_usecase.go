package usecase

import (
	"context"
	"fmt"

	"findMyPhone/internal/domain"
	"findMyPhone/internal/domain/repository"
)

// DeviceUseCase manages device workflows.
type DeviceUseCase struct {
	repo repository.DeviceRepository
}

// NewDeviceUseCase constructs DeviceUseCase.
func NewDeviceUseCase(repo repository.DeviceRepository) *DeviceUseCase {
	return &DeviceUseCase{repo: repo}
}

// CreateDevice validates and persists a device.
func (uc *DeviceUseCase) CreateDevice(ctx context.Context, device *domain.Device) error {
	if device == nil {
		return domain.ErrInvalidInput
	}
	if device.DeviceID == "" || device.IMEI == "" {
		return fmt.Errorf("%w: device_id and imei are required", domain.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, device)
}

// GetDeviceByDeviceID fetches a device by device_id.
func (uc *DeviceUseCase) GetDeviceByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error) {
	if deviceID == "" {
		return nil, domain.ErrInvalidInput
	}
	return uc.repo.GetByDeviceID(ctx, deviceID)
}

// UpdateDevice updates an existing device identified by deviceID.
func (uc *DeviceUseCase) UpdateDevice(ctx context.Context, deviceID string, device *domain.Device) (*domain.Device, error) {
	if device == nil || deviceID == "" {
		return nil, domain.ErrInvalidInput
	}
	if device.IMEI == "" {
		return nil, fmt.Errorf("%w: imei is required", domain.ErrInvalidInput)
	}

	return uc.repo.UpdateByDeviceID(ctx, deviceID, device)
}

// DeleteDevice removes a device by its deviceID.
func (uc *DeviceUseCase) DeleteDevice(ctx context.Context, deviceID string) error {
	if deviceID == "" {
		return domain.ErrInvalidInput
	}

	return uc.repo.DeleteByDeviceID(ctx, deviceID)
}
