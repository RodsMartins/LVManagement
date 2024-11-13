-- name: ListSeedInstructions :many
SELECT * FROM Seed_Instructions;

-- name: GetSeedInstruction :one
SELECT * FROM Seed_Instructions
WHERE seed_instruction_id = $1;

-- name: GetSeedInstructionBySeedId :one
SELECT * FROM Seed_Instructions
WHERE seed_id = $1;

-- name: NewSeedInstruction :one
INSERT INTO Seed_Instructions (seed_instruction_id, seed_id, seed_grams, soaking_hours, stacking_hours, blackout_hours, lights_hours, yield_grams, special_treatment)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateSeedInstruction :exec
UPDATE Seed_Instructions
SET seed_id = $1, seed_grams = $2, soaking_hours = $3, stacking_hours = $4, blackout_hours = $5, lights_hours = $6, yield_grams = $7, special_treatment = $8
WHERE seed_instruction_id = $9;

-- name: DeleteSeedInstruction :execrows
DELETE FROM Seed_Instructions
WHERE seed_instruction_id = $1;
