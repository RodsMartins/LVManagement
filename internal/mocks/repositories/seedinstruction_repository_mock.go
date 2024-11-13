package mocks

import (
	"context"
	"lvm/internal/dtos"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSeedInstructionRepository struct {
    mock.Mock
}

func (m *MockSeedInstructionRepository) DeleteSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (int64, error) {
    args := m.Called(ctx, seedInstructionID)
    return args.Get(0).(int64), args.Error(1)
}

func (m *MockSeedInstructionRepository) GetSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (dtos.SeedInstruction, error) {
    args := m.Called(ctx, seedInstructionID)
    return args.Get(0).(dtos.SeedInstruction), args.Error(1)
}

func (m *MockSeedInstructionRepository) GetSeedInstructionBySeedId(ctx context.Context, seedID uuid.UUID) (dtos.SeedInstruction, error) {
    args := m.Called(ctx, seedID)
    return args.Get(0).(dtos.SeedInstruction), args.Error(1)
}

func (m *MockSeedInstructionRepository) ListSeedInstructions(ctx context.Context) ([]dtos.SeedInstruction, error) {
    args := m.Called(ctx)
    return args.Get(0).([]dtos.SeedInstruction), args.Error(1)
}

func (m *MockSeedInstructionRepository) NewSeedInstruction(ctx context.Context, seedInstruction dtos.SeedInstruction) (dtos.SeedInstruction, error) {
    args := m.Called(ctx, seedInstruction)
    return args.Get(0).(dtos.SeedInstruction), args.Error(1)
}

func (m *MockSeedInstructionRepository) UpdateSeedInstruction(ctx context.Context, seedInstruction dtos.SeedInstruction) error {
    args := m.Called(ctx, seedInstruction)
    return args.Error(0)
}