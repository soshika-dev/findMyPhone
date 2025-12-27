package usecase

import (
	"context"
	"errors"
	"testing"

	"findMyPhone/internal/domain"
)

type mockDeviceRepo struct {
	exists bool
	getErr error
}

func (m *mockDeviceRepo) Create(ctx context.Context, d *domain.Device) error { return nil }
func (m *mockDeviceRepo) GetByDeviceID(ctx context.Context, deviceID string) (*domain.Device, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if !m.exists {
		return nil, domain.ErrNotFound
	}
	return &domain.Device{DeviceID: deviceID}, nil
}

type mockLogRepo struct {
	created   *domain.Log
	createErr error
}

func (m *mockLogRepo) Create(ctx context.Context, l *domain.Log) error {
	m.created = l
	return m.createErr
}
func (m *mockLogRepo) GetLastByDeviceID(ctx context.Context, deviceID string) (*domain.Log, error) {
	return nil, nil
}

func TestCreateLogRequiresDevice(t *testing.T) {
	deviceRepo := &mockDeviceRepo{exists: false}
	logRepo := &mockLogRepo{}
	uc := NewLogUseCase(logRepo, deviceRepo)
	err := uc.CreateLog(context.Background(), &domain.Log{DeviceID: "missing", Latitude: 1, Longitude: 1})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestCreateLogSuccess(t *testing.T) {
	deviceRepo := &mockDeviceRepo{exists: true}
	logRepo := &mockLogRepo{}
	uc := NewLogUseCase(logRepo, deviceRepo)
	entry := &domain.Log{DeviceID: "dev1", Latitude: 1, Longitude: 2}
	if err := uc.CreateLog(context.Background(), entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if logRepo.created != entry {
		t.Fatalf("expected log to be persisted")
	}
}
