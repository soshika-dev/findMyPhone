package repository

import (
	"context"
	"errors"

	"findMyPhone/internal/domain"
	domainRepo "findMyPhone/internal/domain/repository"

	"gorm.io/gorm"
)

// DeviceRepositoryGorm implements DeviceRepository with GORM.
type DeviceRepositoryGorm struct {
	db *gorm.DB
}

// NewDeviceRepository creates a new DeviceRepository.
func NewDeviceRepository(db *gorm.DB) domainRepo.DeviceRepository {
	return &DeviceRepositoryGorm{db: db}
}

// Create inserts a new device.
func (r *DeviceRepositoryGorm) Create(ctx context.Context, device *domain.Device) error {
	if err := r.db.WithContext(ctx).Create(device).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrConflict
		}
		return err
	}
	return nil
}

// GetByDeviceID finds a device by device_id.
func (r *DeviceRepositoryGorm) GetByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error) {
	var device domain.Device
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &device, nil
}
