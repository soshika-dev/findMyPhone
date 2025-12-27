package usecase

import (
	"context"
	"errors"
	"testing"

	"findMyPhone/internal/domain"
)

type mockDeviceRepoUC struct {
	created   *domain.Device
	createErr error
}

func (m *mockDeviceRepoUC) Create(ctx context.Context, d *domain.Device) error {
	m.created = d
	return m.createErr
}

func (m *mockDeviceRepoUC) GetByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error) {
	return nil, nil
}

func TestCreateDeviceValidation(t *testing.T) {
	repo := &mockDeviceRepoUC{}
	uc := NewDeviceUseCase(repo)
	if err := uc.CreateDevice(context.Background(), &domain.Device{}); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestCreateDeviceSuccess(t *testing.T) {
	repo := &mockDeviceRepoUC{}
	uc := NewDeviceUseCase(repo)
	device := &domain.Device{DeviceID: "dev1", IMEI: "imei1"}
	if err := uc.CreateDevice(context.Background(), device); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created != device {
		t.Fatalf("expected device to be persisted")
	}
}
