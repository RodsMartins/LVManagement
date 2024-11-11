package repositories

import (
	"context"
	"lvm/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SeedRepository interface {
	DeleteSeed(ctx context.Context, seedID uuid.UUID) (int64, error)
	GetSeed(ctx context.Context, seedID uuid.UUID) (database.Seed, error)
	ListSeeds(ctx context.Context) ([]database.Seed, error)
	NewSeed(ctx context.Context, params database.NewSeedParams) (database.Seed, error)
	UpdateSeed(ctx context.Context, params database.UpdateSeedParams) error
}

type seedRepository struct {
	queries *database.Queries
}

func NewSeedRepository(queries *database.Queries) SeedRepository {
	return &seedRepository{queries: queries}
}

func (r *seedRepository) DeleteSeed(ctx context.Context, seedID uuid.UUID) (int64, error) {
	return r.queries.DeleteSeed(ctx, pgtype.UUID{Bytes: seedID, Valid: true})
}

func (r *seedRepository) GetSeed(ctx context.Context, seedID uuid.UUID) (database.Seed, error) {
	return r.queries.GetSeed(ctx, pgtype.UUID{Bytes: seedID, Valid: true})
}

func (r *seedRepository) ListSeeds(ctx context.Context) ([]database.Seed, error) {
	return r.queries.ListSeeds(ctx)
}

func (r *seedRepository) NewSeed(ctx context.Context, params database.NewSeedParams) (database.Seed, error) {
	return r.queries.NewSeed(ctx, params)
}

func (r *seedRepository) UpdateSeed(ctx context.Context, params database.UpdateSeedParams) error {
	return r.queries.UpdateSeed(ctx, params)
}