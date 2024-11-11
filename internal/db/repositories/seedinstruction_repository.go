package repositories

import (
	"context"
	"lvm/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SeedInstructionRepository interface {
	DeleteSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (int64, error)
	GetSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (database.SeedInstruction, error)
	GetSeedInstructionsBySeedId(ctx context.Context, seedID uuid.UUID) (database.SeedInstruction, error)
	ListSeedInstructions(ctx context.Context) ([]database.SeedInstruction, error)
	NewSeedInstruction(ctx context.Context, params database.NewSeedInstructionParams) (database.SeedInstruction, error)
	UpdateSeedInstruction(ctx context.Context, params database.UpdateSeedInstructionParams) error
}

type seedInstructionRepository struct {
	queries *database.Queries
}

func NewSeedInstructionRepository(queries *database.Queries) SeedInstructionRepository {
	return &seedInstructionRepository{queries: queries}
}

func (r *seedInstructionRepository) DeleteSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (int64, error) {
	return r.queries.DeleteSeedInstruction(ctx, pgtype.UUID{Bytes: seedInstructionID})
}

func (r *seedInstructionRepository) GetSeedInstruction(ctx context.Context, seedInstructionID uuid.UUID) (database.SeedInstruction, error) {
	return r.queries.GetSeedInstruction(ctx, pgtype.UUID{Bytes: seedInstructionID})
}

func (r *seedInstructionRepository) GetSeedInstructionsBySeedId(ctx context.Context, seedID uuid.UUID) (database.SeedInstruction, error) {
	return r.queries.GetSeedInstructionsBySeedId(ctx, pgtype.UUID{Bytes: seedID})
}

func (r *seedInstructionRepository) ListSeedInstructions(ctx context.Context) ([]database.SeedInstruction, error) {
	return r.queries.ListSeedInstructions(ctx)
}

func (r *seedInstructionRepository) NewSeedInstruction(ctx context.Context, params database.NewSeedInstructionParams) (database.SeedInstruction, error) {
	return r.queries.NewSeedInstruction(ctx, params)
}

func (r *seedInstructionRepository) UpdateSeedInstruction(ctx context.Context, params database.UpdateSeedInstructionParams) error {
	return r.queries.UpdateSeedInstruction(ctx, params)
}