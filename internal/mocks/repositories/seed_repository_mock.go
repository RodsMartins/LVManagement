package mocks

import (
	"context"
	"lvm/internal/dtos"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSeedRepository struct {
    mock.Mock
}

func (m *MockSeedRepository) DeleteSeed(ctx context.Context, seedID uuid.UUID) (int64, error) {
    args := m.Called(ctx, seedID)
    return args.Get(0).(int64), args.Error(1)
}

func (m *MockSeedRepository) GetSeed(ctx context.Context, seedID uuid.UUID) (dtos.Seed, error) {
    args := m.Called(ctx, seedID)
    return args.Get(0).(dtos.Seed), args.Error(1)
}

func (m *MockSeedRepository) ListSeeds(ctx context.Context) ([]dtos.Seed, error) {
    args := m.Called(ctx)
    return args.Get(0).([]dtos.Seed), args.Error(1)
}

func (m *MockSeedRepository) NewSeed(ctx context.Context, seed dtos.Seed) (dtos.Seed, error) {
    args := m.Called(ctx, seed)
    return args.Get(0).(dtos.Seed), args.Error(1)
}

func (m *MockSeedRepository) UpdateSeed(ctx context.Context, seed dtos.Seed) error {
    args := m.Called(ctx, seed)
    return args.Error(0)
}