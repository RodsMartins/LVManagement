package repositories

import (
	"context"
	"fmt"
	"lvm/database"
	"lvm/internal/dtos"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CropRepository interface {
	GetCrop(ctx context.Context, cropID uuid.UUID) (dtos.Crop, error)
	GetExistingCodesForSeed(ctx context.Context, input GetExistingCodesForSeedInput) ([]string, error)
	ListCrop(ctx context.Context) ([]dtos.Crop, error)
	NewCrop(ctx context.Context, params database.NewCropParams) (dtos.Crop, error)
}

type cropRepository struct {
	queries *database.Queries
}

type GetExistingCodesForSeedInput struct {
	FirstLetter string
	SeedID      uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
}

func NewCropRepository(queries *database.Queries) CropRepository {
	return &cropRepository{queries: queries}
}

func (r *cropRepository) GetCrop(ctx context.Context, cropID uuid.UUID) (dtos.Crop, error) {
	newCrop, err := r.queries.GetCrop(ctx, pgtype.UUID{Bytes: cropID, Valid: true})
	if err != nil {
		return dtos.Crop{}, err
	}

	return dtos.CropFromDatabaseModel(newCrop), nil
}

func (r *cropRepository) GetExistingCodesForSeed(ctx context.Context, input GetExistingCodesForSeedInput) ([]string, error) {
	params := database.GetExistingCodesForSeedParams{
		FirstLetter: pgtype.Text{String: input.FirstLetter, Valid: true},
		StartDate:   pgtype.Timestamp{Time: input.StartDate, Valid: true},
		EndDate:     pgtype.Timestamp{Time: input.EndDate, Valid: true},
		SeedID:      pgtype.UUID{Bytes: input.SeedID, Valid: true},
	}

	codes, err := r.queries.GetExistingCodesForSeed(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("getting existing codes: %w", err)
	}

	// Convert pgtype.Text to []string
	result := make([]string, len(codes))
	for i, code := range codes {
		if code.Valid {
			result[i] = code.String
		}
	}

	return result, nil
}

func (r *cropRepository) ListCrop(ctx context.Context) ([]dtos.Crop, error) {
	crops, err := r.queries.ListCrop(ctx)
	if err != nil {
		return nil, err
	}

	dtoCrops := make([]dtos.Crop, len(crops))
	for i, c := range crops {
		dtoCrops[i] = dtos.CropFromDatabaseModel(c)
	}

	return dtoCrops, nil
}

func (r *cropRepository) NewCrop(ctx context.Context, params database.NewCropParams) (dtos.Crop, error) {
	crop, err := r.queries.NewCrop(ctx, params)
	if err != nil {
		return dtos.Crop{}, err
	}

	return dtos.CropFromDatabaseModel(crop), nil
}
