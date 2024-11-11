-- name: ListSeeds :many
SELECT * FROM Seeds;

-- name: GetSeed :one
SELECT * FROM Seeds
WHERE seed_id = $1;

-- name: NewSeed :one
INSERT INTO Seeds (seed_id, name, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSeed :exec
UPDATE Seeds
SET name = $1, type = $2
WHERE seed_id = $3;

-- name: DeleteSeed :execrows
DELETE FROM Seeds
WHERE seed_id = $1;
