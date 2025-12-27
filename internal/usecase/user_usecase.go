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
