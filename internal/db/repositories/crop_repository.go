package repositories

import (
	"context"
	"lvm/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CropRepository interface {
	GetCrop(ctx context.Context, cropID uuid.UUID) (database.Crop, error)
	GetExistingCodesForSeed(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error)
	ListCrop(ctx context.Context) ([]database.Crop, error)
	NewCrop(ctx context.Context, params database.NewCropParams) (database.Crop, error)
}

type cropRepository struct {
	queries *database.Queries
}

func NewCropRepository(queries *database.Queries) CropRepository {
	return &cropRepository{queries: queries}
}

func (r *cropRepository) GetCrop(ctx context.Context, cropID uuid.UUID) (database.Crop, error) {
	return r.queries.GetCrop(ctx, pgtype.UUID{Bytes: cropID, Valid: true})
}

func (r *cropRepository) GetExistingCodesForSeed(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error) {
	return r.queries.GetExistingCodesForSeed(ctx, params)
}

func (r *cropRepository) ListCrop(ctx context.Context) ([]database.Crop, error) {
	var crops []database.Crop
	crop, err := r.queries.ListCrop(ctx)
	if err != nil {
		return nil, err
	}
	crops = append(crops, crop)
	return crops, nil
}

func (r *cropRepository) NewCrop(ctx context.Context, params database.NewCropParams) (database.Crop, error) {
	return r.queries.NewCrop(ctx, params)
}