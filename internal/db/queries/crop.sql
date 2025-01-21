-- name: ListCrop :many
SELECT * FROM Crops;

-- name: GetCrop :one
SELECT *
FROM Crops
WHERE crop_id = $1;

-- name: NewCrop :one
INSERT INTO Crops (crop_id, seed_id, soaking_start, stacking_start, blackout_start, lights_start, harvest, code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetExistingCodesForSeed :many
-- @arg first_letter string
WITH reference_letter AS (
    SELECT LEFT(name, 1) as first_letter
    FROM Seeds s
    WHERE s.seed_id = sqlc.arg('seed_id')
)
SELECT c.code
FROM Crops c
WHERE c.code ~ ('^' || sqlc.arg('first_letter') || '\d+$')
  AND (
    c.stacking_start BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date')
    OR (c.soaking_start IS NOT NULL AND c.soaking_start BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date'))
    OR (c.harvest BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date'))
  )
ORDER BY c.code;

-- name: ListCropsByDate :many
-- @arg date
SELECT 
  *
FROM Crops c
WHERE 
    (c.soaking_start IS NOT NULL AND sqlc.arg('date') BETWEEN c.soaking_start AND c.harvest)
    OR (sqlc.arg('date') BETWEEN c.stacking_start AND c.harvest);