// Code generated by sqlc. DO NOT EDIT.
// source: messages.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (
    content, external_id
) VALUES (
    $1, $2
) RETURNING id, created_at, content, external_id
`

type CreateMessageParams struct {
	Content    string
	ExternalID []uuid.UUID
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage, arg.Content, pq.Array(arg.ExternalID))
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Content,
		pq.Array(&i.ExternalID),
	)
	return i, err
}

const getMessage = `-- name: GetMessage :one
SELECT id, created_at, content, external_id
FROM messages
WHERE id = $1
`

func (q *Queries) GetMessage(ctx context.Context, id uuid.UUID) (Message, error) {
	row := q.db.QueryRowContext(ctx, getMessage, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Content,
		pq.Array(&i.ExternalID),
	)
	return i, err
}

const listMessages = `-- name: ListMessages :many
SELECT id, created_at, content, external_id
FROM messages
ORDER BY created_at
`

func (q *Queries) ListMessages(ctx context.Context) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, listMessages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Content,
			pq.Array(&i.ExternalID),
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMessage = `-- name: UpdateMessage :exec
UPDATE messages
SET content = $2, external_id = $3
WHERE id = $1
`

type UpdateMessageParams struct {
	ID         uuid.UUID
	Content    string
	ExternalID []uuid.UUID
}

func (q *Queries) UpdateMessage(ctx context.Context, arg UpdateMessageParams) error {
	_, err := q.db.ExecContext(ctx, updateMessage, arg.ID, arg.Content, pq.Array(arg.ExternalID))
	return err
}
