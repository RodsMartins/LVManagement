package repositories

import (
	"context"
	"lvm/database"
	"lvm/internal/dtos"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SeedRepository interface {
	DeleteSeed(ctx context.Context, seedID uuid.UUID) (int64, error)
	GetSeed(ctx context.Context, seedID uuid.UUID) (dtos.Seed, error)
	ListSeeds(ctx context.Context) ([]dtos.Seed, error)
	NewSeed(ctx context.Context, seed dtos.Seed) (dtos.Seed, error)
	UpdateSeed(ctx context.Context, seed dtos.Seed) error
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

func (r *seedRepository) GetSeed(ctx context.Context, seedID uuid.UUID) (dtos.Seed, error) {
	seed, err := r.queries.GetSeed(ctx, pgtype.UUID{Bytes: seedID, Valid: true})
	if err != nil {
		return dtos.Seed{}, err
	}

	return dtos.Seed{}.FromDatabaseModel(seed), nil
}

func (r *seedRepository) ListSeeds(ctx context.Context) ([]dtos.Seed, error) {
	seeds, err := r.queries.ListSeeds(ctx)
	if err != nil {
		return nil, err
	}

	dtosSeed := make([]dtos.Seed, len(seeds))
	for i, seed := range seeds {
		dtosSeed[i] = dtos.Seed{}.FromDatabaseModel(seed)
	}

	return dtosSeed, nil
}

func (r *seedRepository) NewSeed(ctx context.Context, seed dtos.Seed) (dtos.Seed, error) {
	dbSeed, err := r.queries.NewSeed(ctx, seed.ToNewSeedParams())
	if err != nil {
		return dtos.Seed{}, err
	}

	return dtos.Seed{}.FromDatabaseModel(dbSeed), nil
}

func (r *seedRepository) UpdateSeed(ctx context.Context, seed dtos.Seed) error {
	return r.queries.UpdateSeed(ctx, seed.ToUpdateSeedParams())
}