// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: seedinstructions.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteSeedInstruction = `-- name: DeleteSeedInstruction :execrows
DELETE FROM Seed_Instructions
WHERE seed_instruction_id = $1
`

func (q *Queries) DeleteSeedInstruction(ctx context.Context, seedInstructionID pgtype.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteSeedInstruction, seedInstructionID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getSeedInstruction = `-- name: GetSeedInstruction :one
SELECT seed_instruction_id, seed_id, seed_grams, soaking_hours, stacking_hours, blackout_hours, lights_hours, yield_grams, special_treatment FROM Seed_Instructions
WHERE seed_instruction_id = $1
`

func (q *Queries) GetSeedInstruction(ctx context.Context, seedInstructionID pgtype.UUID) (SeedInstruction, error) {
	row := q.db.QueryRow(ctx, getSeedInstruction, seedInstructionID)
	var i SeedInstruction
	err := row.Scan(
		&i.SeedInstructionID,
		&i.SeedID,
		&i.SeedGrams,
		&i.SoakingHours,
		&i.StackingHours,
		&i.BlackoutHours,
		&i.LightsHours,
		&i.YieldGrams,
		&i.SpecialTreatment,
	)
	return i, err
}

const getSeedInstructionsBySeedId = `-- name: GetSeedInstructionsBySeedId :one
SELECT seed_instruction_id, seed_id, seed_grams, soaking_hours, stacking_hours, blackout_hours, lights_hours, yield_grams, special_treatment FROM Seed_Instructions
WHERE seed_id = $1
`

func (q *Queries) GetSeedInstructionsBySeedId(ctx context.Context, seedID pgtype.UUID) (SeedInstruction, error) {
	row := q.db.QueryRow(ctx, getSeedInstructionsBySeedId, seedID)
	var i SeedInstruction
	err := row.Scan(
		&i.SeedInstructionID,
		&i.SeedID,
		&i.SeedGrams,
		&i.SoakingHours,
		&i.StackingHours,
		&i.BlackoutHours,
		&i.LightsHours,
		&i.YieldGrams,
		&i.SpecialTreatment,
	)
	return i, err
}

const listSeedInstructions = `-- name: ListSeedInstructions :many
SELECT seed_instruction_id, seed_id, seed_grams, soaking_hours, stacking_hours, blackout_hours, lights_hours, yield_grams, special_treatment FROM Seed_Instructions
`

func (q *Queries) ListSeedInstructions(ctx context.Context) ([]SeedInstruction, error) {
	rows, err := q.db.Query(ctx, listSeedInstructions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SeedInstruction
	for rows.Next() {
		var i SeedInstruction
		if err := rows.Scan(
			&i.SeedInstructionID,
			&i.SeedID,
			&i.SeedGrams,
			&i.SoakingHours,
			&i.StackingHours,
			&i.BlackoutHours,
			&i.LightsHours,
			&i.YieldGrams,
			&i.SpecialTreatment,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newSeedInstruction = `-- name: NewSeedInstruction :one
INSERT INTO Seed_Instructions (seed_instruction_id, seed_id, seed_grams, soaking_hours, stacking_hours, blackout_hours, lights_hours, yield_grams, special_treatment)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING seed_instruction_id, seed_id, seed_grams, soaking_hours, stacking_hours, blackout_hours, lights_hours, yield_grams, special_treatment
`

type NewSeedInstructionParams struct {
	SeedInstructionID pgtype.UUID
	SeedID            pgtype.UUID
	SeedGrams         pgtype.Int4
	SoakingHours      pgtype.Int4
	StackingHours     pgtype.Int4
	BlackoutHours     pgtype.Int4
	LightsHours       pgtype.Int4
	YieldGrams        pgtype.Int4
	SpecialTreatment  pgtype.Text
}

func (q *Queries) NewSeedInstruction(ctx context.Context, arg NewSeedInstructionParams) (SeedInstruction, error) {
	row := q.db.QueryRow(ctx, newSeedInstruction,
		arg.SeedInstructionID,
		arg.SeedID,
		arg.SeedGrams,
		arg.SoakingHours,
		arg.StackingHours,
		arg.BlackoutHours,
		arg.LightsHours,
		arg.YieldGrams,
		arg.SpecialTreatment,
	)
	var i SeedInstruction
	err := row.Scan(
		&i.SeedInstructionID,
		&i.SeedID,
		&i.SeedGrams,
		&i.SoakingHours,
		&i.StackingHours,
		&i.BlackoutHours,
		&i.LightsHours,
		&i.YieldGrams,
		&i.SpecialTreatment,
	)
	return i, err
}

const updateSeedInstruction = `-- name: UpdateSeedInstruction :exec
UPDATE Seed_Instructions
SET seed_id = $1, seed_grams = $2, soaking_hours = $3, stacking_hours = $4, blackout_hours = $5, lights_hours = $6, yield_grams = $7, special_treatment = $8
WHERE seed_instruction_id = $9
`

type UpdateSeedInstructionParams struct {
	SeedID            pgtype.UUID
	SeedGrams         pgtype.Int4
	SoakingHours      pgtype.Int4
	StackingHours     pgtype.Int4
	BlackoutHours     pgtype.Int4
	LightsHours       pgtype.Int4
	YieldGrams        pgtype.Int4
	SpecialTreatment  pgtype.Text
	SeedInstructionID pgtype.UUID
}

func (q *Queries) UpdateSeedInstruction(ctx context.Context, arg UpdateSeedInstructionParams) error {
	_, err := q.db.Exec(ctx, updateSeedInstruction,
		arg.SeedID,
		arg.SeedGrams,
		arg.SoakingHours,
		arg.StackingHours,
		arg.BlackoutHours,
		arg.LightsHours,
		arg.YieldGrams,
		arg.SpecialTreatment,
		arg.SeedInstructionID,
	)
	return err
}