-- name: List :one
SELECT * FROM Crops;

-- name: Get :one
SELECT *
FROM Crops
WHERE crop_id = $1;