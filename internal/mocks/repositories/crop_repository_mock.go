package mocks

import (
	"context"
	"lvm/database"
	"lvm/internal/db/repositories"
	"lvm/internal/dtos"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCropRepository struct {
    mock.Mock
}

func (m *MockCropRepository) GetCrop(ctx context.Context, cropID uuid.UUID) (dtos.Crop, error) {
    args := m.Called(ctx, cropID)
    return args.Get(0).(dtos.Crop), args.Error(1)
}

func (m *MockCropRepository) GetExistingCodesForSeed(ctx context.Context, input repositories.GetExistingCodesForSeedInput) ([]string, error) {
    args := m.Called(ctx, input)
    return args.Get(0).([]string), args.Error(1)
}

func (m *MockCropRepository) ListCrop(ctx context.Context) ([]dtos.Crop, error) {
    args := m.Called(ctx)
    return args.Get(0).([]dtos.Crop), args.Error(1)
}

func (m *MockCropRepository) NewCrop(ctx context.Context, params database.NewCropParams) (dtos.Crop, error) {
    args := m.Called(ctx, params)
    return args.Get(0).(dtos.Crop), args.Error(1)
}