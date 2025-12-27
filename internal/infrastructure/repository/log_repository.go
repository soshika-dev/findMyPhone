package repository

import (
	"context"

	"findMyPhone/internal/domain"
	domainRepo "findMyPhone/internal/domain/repository"

	"gorm.io/gorm"
)

// LogRepositoryGorm implements log persistence.
type LogRepositoryGorm struct {
	db *gorm.DB
}

// NewLogRepository creates a new LogRepository.
func NewLogRepository(db *gorm.DB) domainRepo.LogRepository {
	return &LogRepositoryGorm{db: db}
}

// Create inserts a new log entry.
func (r *LogRepositoryGorm) Create(ctx context.Context, log *domain.Log) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return err
	}
	return nil
}

// GetLastByDeviceID retrieves most recent log for device.
func (r *LogRepositoryGorm) GetLastByDeviceID(ctx context.Context, deviceID string) (*domain.Log, error) {
	var l domain.Log
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).
		Order("created_at desc").First(&l).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}
