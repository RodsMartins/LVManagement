-- name: CreateTmp :one
INSERT INTO Fertilizers (fertilizer_id, name)
VALUES ($1, $2)
RETURNING *;
