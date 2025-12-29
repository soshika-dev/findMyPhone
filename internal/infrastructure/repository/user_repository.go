package repository

import (
	"context"

	"findMyPhone/internal/domain"
	domainRepo "findMyPhone/internal/domain/repository"

	"gorm.io/gorm"
)

// UserRepositoryGorm implements UserRepository using GORM.
type UserRepositoryGorm struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepositoryGorm.
func NewUserRepository(db *gorm.DB) domainRepo.UserRepository {
	return &UserRepositoryGorm{db: db}
}

// Create inserts a new user.
func (r *UserRepositoryGorm) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateError(err) {
			return domain.ErrConflict
		}
		return err
	}
	return nil
}

// GetByDeviceID fetches a user by device identifier.
func (r *UserRepositoryGorm) GetByDeviceID(ctx context.Context, deviceID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateByDeviceID updates a user tied to the provided device identifier.
func (r *UserRepositoryGorm) UpdateByDeviceID(ctx context.Context, deviceID string, user *domain.User) (*domain.User, error) {
	var existing domain.User
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"name":         user.Name,
		"phone":        user.Phone,
		"backup_phone": user.BackupPhone,
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
