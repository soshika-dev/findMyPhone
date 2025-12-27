package repository

import (
	"context"
	"errors"

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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
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
