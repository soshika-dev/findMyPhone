package usecase

import (
	"context"
	"fmt"

	"findMyPhone/internal/domain"
	"findMyPhone/internal/domain/repository"
)

// LogUseCase handles log workflows.
type LogUseCase struct {
	logRepo    repository.LogRepository
	deviceRepo repository.DeviceRepository
}

// NewLogUseCase constructs a LogUseCase.
func NewLogUseCase(logRepo repository.LogRepository, deviceRepo repository.DeviceRepository) *LogUseCase {
	return &LogUseCase{logRepo: logRepo, deviceRepo: deviceRepo}
}

// CreateLog validates and stores a log ensuring device exists.
func (uc *LogUseCase) CreateLog(ctx context.Context, log *domain.Log) error {
	if log == nil {
		return domain.ErrInvalidInput
	}
	if log.DeviceID == "" {
		return fmt.Errorf("%w: device_id is required", domain.ErrInvalidInput)
	}
	_, err := uc.deviceRepo.GetByDeviceID(ctx, log.DeviceID)
	if err != nil {
		return err
	}
	return uc.logRepo.Create(ctx, log)
}

// GetLastLogByDeviceID retrieves latest log for device.
func (uc *LogUseCase) GetLastLogByDeviceID(ctx context.Context, deviceID string) (*domain.Log, error) {
	if deviceID == "" {
		return nil, domain.ErrInvalidInput
	}
	return uc.logRepo.GetLastByDeviceID(ctx, deviceID)
}
