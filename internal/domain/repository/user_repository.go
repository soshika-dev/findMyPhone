package repository

import (
	"context"

	"findMyPhone/internal/domain"
)

// UserRepository defines persistence for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByDeviceID(ctx context.Context, deviceID string) (*domain.User, error)
}
