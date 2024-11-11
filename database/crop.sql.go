// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: crop.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getCrop = `-- name: GetCrop :one
SELECT crop_id, seed_id, soaking_start, stacking_start, blackout_start, lights_start, harvest, code, yield_grams
FROM Crops
WHERE crop_id = $1
`

func (q *Queries) GetCrop(ctx context.Context, cropID pgtype.UUID) (Crop, error) {
	row := q.db.QueryRow(ctx, getCrop, cropID)
	var i Crop
	err := row.Scan(
		&i.CropID,
		&i.SeedID,
		&i.SoakingStart,
		&i.StackingStart,
		&i.BlackoutStart,
		&i.LightsStart,
		&i.Harvest,
		&i.Code,
		&i.YieldGrams,
	)
	return i, err
}

const getExistingCodesForSeed = `-- name: GetExistingCodesForSeed :many
WITH reference_letter AS (
    SELECT LEFT(name, 1) as first_letter
    FROM Seeds s
    WHERE s.seed_id = $4
)
SELECT c.code
FROM Crops c
WHERE c.code ~ ('^' || $1 || '\d+$')
  AND (
    c.stacking_start BETWEEN $2 AND $3
    OR (c.soaking_start IS NOT NULL AND c.soaking_start BETWEEN $2 AND $3)
    OR (c.harvest BETWEEN $2 AND $3)
  )
ORDER BY c.code
`

type GetExistingCodesForSeedParams struct {
	FirstLetter pgtype.Text
	StartDate   pgtype.Timestamp
	EndDate     pgtype.Timestamp
	SeedID      pgtype.UUID
}

// @arg first_letter string
func (q *Queries) GetExistingCodesForSeed(ctx context.Context, arg GetExistingCodesForSeedParams) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, getExistingCodesForSeed,
		arg.FirstLetter,
		arg.StartDate,
		arg.EndDate,
		arg.SeedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var code pgtype.Text
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		items = append(items, code)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCrop = `-- name: ListCrop :one
SELECT crop_id, seed_id, soaking_start, stacking_start, blackout_start, lights_start, harvest, code, yield_grams FROM Crops
`

func (q *Queries) ListCrop(ctx context.Context) (Crop, error) {
	row := q.db.QueryRow(ctx, listCrop)
	var i Crop
	err := row.Scan(
		&i.CropID,
		&i.SeedID,
		&i.SoakingStart,
		&i.StackingStart,
		&i.BlackoutStart,
		&i.LightsStart,
		&i.Harvest,
		&i.Code,
		&i.YieldGrams,
	)
	return i, err
}

const newCrop = `-- name: NewCrop :one
INSERT INTO Crops (crop_id, seed_id, soaking_start, stacking_start, blackout_start, lights_start, harvest, code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING crop_id, seed_id, soaking_start, stacking_start, blackout_start, lights_start, harvest, code, yield_grams
`

type NewCropParams struct {
	CropID        pgtype.UUID
	SeedID        pgtype.UUID
	SoakingStart  pgtype.Timestamp
	StackingStart pgtype.Timestamp
	BlackoutStart pgtype.Timestamp
	LightsStart   pgtype.Timestamp
	Harvest       pgtype.Timestamp
	Code          pgtype.Text
}

func (q *Queries) NewCrop(ctx context.Context, arg NewCropParams) (Crop, error) {
	row := q.db.QueryRow(ctx, newCrop,
		arg.CropID,
		arg.SeedID,
		arg.SoakingStart,
		arg.StackingStart,
		arg.BlackoutStart,
		arg.LightsStart,
		arg.Harvest,
		arg.Code,
	)
	var i Crop
	err := row.Scan(
		&i.CropID,
		&i.SeedID,
		&i.SoakingStart,
		&i.StackingStart,
		&i.BlackoutStart,
		&i.LightsStart,
		&i.Harvest,
		&i.Code,
		&i.YieldGrams,
	)
	return i, err
}
