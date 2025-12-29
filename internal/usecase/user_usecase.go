package usecase

import (
	"context"
	"fmt"

	"findMyPhone/internal/domain"
	"findMyPhone/internal/domain/repository"
)

// UserUseCase handles business logic for users.
type UserUseCase struct {
	repo repository.UserRepository
}

// NewUserUseCase creates a new UserUseCase instance.
func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

// CreateUser stores a new user after validating required fields.
func (uc *UserUseCase) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return domain.ErrInvalidInput
	}
	if user.Name == "" || user.DeviceID == "" || user.Phone == "" {
		return fmt.Errorf("%w: name, device_id, and phone are required", domain.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, user)
}

// GetUserByDeviceID retrieves a user by their device id.
func (uc *UserUseCase) GetUserByDeviceID(ctx context.Context, deviceID string) (*domain.User, error) {
	if deviceID == "" {
		return nil, domain.ErrInvalidInput
	}
	return uc.repo.GetByDeviceID(ctx, deviceID)
}

// UpdateUser updates fields for the user associated with the deviceID.
func (uc *UserUseCase) UpdateUser(ctx context.Context, deviceID string, user *domain.User) (*domain.User, error) {
	if user == nil || deviceID == "" {
		return nil, domain.ErrInvalidInput
	}
	if user.Name == "" || user.Phone == "" {
		return nil, fmt.Errorf("%w: name and phone are required", domain.ErrInvalidInput)
	}

	return uc.repo.UpdateByDeviceID(ctx, deviceID, user)
}
