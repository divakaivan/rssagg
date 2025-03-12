// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: logs.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createLog = `-- name: CreateLog :exec
insert into logs (
    timestamp,
    caller_user_id,
    method,
    url,
    status,
    duration_ms
) values ($1, $2, $3, $4, $5, $6)
`

type CreateLogParams struct {
	Timestamp    time.Time
	CallerUserID uuid.UUID
	Method       string
	Url          string
	Status       string
	DurationMs   int64
}

func (q *Queries) CreateLog(ctx context.Context, arg CreateLogParams) error {
	_, err := q.db.ExecContext(ctx, createLog,
		arg.Timestamp,
		arg.CallerUserID,
		arg.Method,
		arg.Url,
		arg.Status,
		arg.DurationMs,
	)
	return err
}
