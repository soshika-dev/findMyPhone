package usecase

import (
	"context"
	"errors"
	"testing"

	"findMyPhone/internal/domain"
)

type mockUserRepo struct {
	created   *domain.User
	getResp   *domain.User
	getErr    error
	createErr error
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	m.created = user
	return m.createErr
}

func (m *mockUserRepo) GetByDeviceID(ctx context.Context, deviceID string) (*domain.User, error) {
	return m.getResp, m.getErr
}

func TestCreateUserValidation(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewUserUseCase(repo)
	err := uc.CreateUser(context.Background(), &domain.User{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}

func TestCreateUserSuccess(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewUserUseCase(repo)
	user := &domain.User{Name: "John", DeviceID: "dev1", Phone: "123"}
	if err := uc.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created != user {
		t.Fatalf("expected user to be persisted")
	}
}
