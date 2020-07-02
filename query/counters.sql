-- name: ListCounters :many
SELECT *
FROM counters;

-- name: GetCounter :one
SELECT *
FROM counters
WHERE id = $1;

-- name: CreateCounter :one
INSERT INTO counters (
    value
) VALUES (
    $1
) RETURNING *;

-- name: UpdateCounter :exec
UPDATE counters
SET value = $2
WHERE id = $1;
