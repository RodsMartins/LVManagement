-- name: ListSeeds :many
SELECT * FROM Seeds;

-- name: NewSeed :one
INSERT INTO Seeds (seed_id, name, type)
VALUES ($1, $2, $3)
RETURNING *;
