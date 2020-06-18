-- name: ListMessages :many
SELECT *
FROM messages
ORDER BY created_at;

-- name: GetMessage :one
SELECT *
FROM messages
WHERE id = $1;

-- name: CreateMessage :one
INSERT INTO messages (
    content, external_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: UpdateMessage :exec
UPDATE messages
SET content = $2, external_id = $3
WHERE id = $1;
