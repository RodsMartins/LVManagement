package repositories

import (
	"context"
	"lvm/database"
	"lvm/internal/dtos"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SeedInstructionRepository interface {
	DeleteSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (int64, error)
	GetSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (dtos.SeedInstruction, error)
	GetSeedInstructionBySeedId(ctx context.Context, seedID uuid.UUID) (dtos.SeedInstruction, error)
	ListSeedInstructions(ctx context.Context) ([]dtos.SeedInstruction, error)
	NewSeedInstruction(ctx context.Context, seedInstruction dtos.SeedInstruction) (dtos.SeedInstruction, error)
	UpdateSeedInstruction(ctx context.Context, seedInstruction dtos.SeedInstruction) error
}

type seedInstructionRepository struct {
	queries *database.Queries
}

type SeedInstructionParamsInput struct {
	SeedInstructionID uuid.UUID
	SeedID            uuid.UUID
	SeedGrams         int
	SoakingHours      int
	StackingHours     int
	BlackoutHours     int
	LightsHours       int
	YieldGrams        int
	SpecialTreatment  string
}

func NewSeedInstructionRepository(queries *database.Queries) SeedInstructionRepository {
	return &seedInstructionRepository{queries: queries}
}

func (r *seedInstructionRepository) DeleteSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (int64, error) {
	return r.queries.DeleteSeedInstruction(ctx, pgtype.UUID{Bytes: seedInstructionID})
}

func (r *seedInstructionRepository) GetSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (dtos.SeedInstruction, error) {
	seedinstruction, err := r.queries.GetSeedInstruction(ctx, pgtype.UUID{Bytes: seedInstructionID})
	if err != nil {
		return dtos.SeedInstruction{}, err
	}

	return dtos.SeedInstruction{}.FromDatabaseModel(seedinstruction), nil
}

func (r *seedInstructionRepository) GetSeedInstructionBySeedId(ctx context.Context, seedID uuid.UUID) (dtos.SeedInstruction, error) {
	seedinstruction, err :=  r.queries.GetSeedInstructionBySeedId(ctx, pgtype.UUID{Bytes: seedID, Valid: true})
	if err != nil {
		return dtos.SeedInstruction{}, err
	}

	return dtos.SeedInstruction{}.FromDatabaseModel(seedinstruction), nil
}

func (r *seedInstructionRepository) ListSeedInstructions(ctx context.Context) ([]dtos.SeedInstruction, error) {
	seedInstructions, err := r.queries.ListSeedInstructions(ctx)
	if err != nil {
		return nil, err
	}

	dtosSeedInstructions := make([]dtos.SeedInstruction, len(seedInstructions))
	for i, seedInstruction := range seedInstructions {
		dtosSeedInstructions[i] = dtos.SeedInstruction{}.FromDatabaseModel(seedInstruction)
	}

	return dtosSeedInstructions, nil
}

func (r *seedInstructionRepository) NewSeedInstruction(ctx context.Context, seedInstruction dtos.SeedInstruction) (dtos.SeedInstruction, error) {
	dbSeedInstruction, err := r.queries.NewSeedInstruction(ctx, seedInstruction.ToNewSeedInstructionParams())
	if err != nil {
		return dtos.SeedInstruction{}, err
	}

	return dtos.SeedInstruction{}.FromDatabaseModel(dbSeedInstruction), nil
}

func (r *seedInstructionRepository) UpdateSeedInstruction(ctx context.Context, seedInstruction dtos.SeedInstruction) error {
	return r.queries.UpdateSeedInstruction(ctx, seedInstruction.ToUpdateSeedInstructionParams())
}