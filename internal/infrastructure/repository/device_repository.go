package repository

import (
	"context"

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
		if isDuplicateError(err) {
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

// UpdateByDeviceID updates a device matching the provided deviceID.
func (r *DeviceRepositoryGorm) UpdateByDeviceID(ctx context.Context, deviceID string, device *domain.Device) (*domain.Device, error) {
	var existing domain.Device
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"imei":       device.IMEI,
		"generation": device.Generation,
		"name":       device.Name,
		"lost":       device.Lost,
	}

	if err := r.db.WithContext(ctx).Model(&existing).Updates(updates).Error; err != nil {
		if isDuplicateError(err) {
			return nil, domain.ErrConflict
		}
		return nil, err
	}

	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&existing).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

// DeleteByDeviceID removes a device matching the given deviceID.
func (r *DeviceRepositoryGorm) DeleteByDeviceID(ctx context.Context, deviceID string) error {
	var existing domain.Device
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&existing).Error; err != nil {
		return err
	}

	return nil
}
