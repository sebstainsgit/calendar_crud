// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: events.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (event_id, event_name, author_id, date, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING event_id, event_name, author_id, date, created_at, updated_at
`

type CreateEventParams struct {
	EventID   uuid.UUID
	EventName string
	AuthorID  uuid.UUID
	Date      time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, createEvent,
		arg.EventID,
		arg.EventName,
		arg.AuthorID,
		arg.Date,
		arg.UpdatedAt,
		arg.CreatedAt,
	)
	var i Event
	err := row.Scan(
		&i.EventID,
		&i.EventName,
		&i.AuthorID,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM events
WHERE event_id = $1
`

func (q *Queries) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEvent, eventID)
	return err
}

const getEventByID = `-- name: GetEventByID :one
SELECT event_id, event_name, author_id, date, created_at, updated_at FROM events
WHERE event_id = $1
`

func (q *Queries) GetEventByID(ctx context.Context, eventID uuid.UUID) (Event, error) {
	row := q.db.QueryRowContext(ctx, getEventByID, eventID)
	var i Event
	err := row.Scan(
		&i.EventID,
		&i.EventName,
		&i.AuthorID,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getEventByName = `-- name: GetEventByName :one
SELECT event_id, event_name, author_id, date, created_at, updated_at FROM events
WHERE event_name = $1
`

func (q *Queries) GetEventByName(ctx context.Context, eventName string) (Event, error) {
	row := q.db.QueryRowContext(ctx, getEventByName, eventName)
	var i Event
	err := row.Scan(
		&i.EventID,
		&i.EventName,
		&i.AuthorID,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getEvents = `-- name: GetEvents :many
SELECT event_id, event_name, author_id, date, created_at, updated_at FROM events
ORDER BY date DESC
`

func (q *Queries) GetEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.EventID,
			&i.EventName,
			&i.AuthorID,
			&i.Date,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getUserMadeEvents = `-- name: GetUserMadeEvents :many
SELECT event_id, event_name, author_id, date, created_at, updated_at FROM events
WHERE author_id = $1
`

func (q *Queries) GetUserMadeEvents(ctx context.Context, authorID uuid.UUID) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getUserMadeEvents, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.EventID,
			&i.EventName,
			&i.AuthorID,
			&i.Date,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateEvent = `-- name: UpdateEvent :one
UPDATE events
SET date = $2, event_name = $3, updated_at = $4
WHERE event_id = $1
RETURNING event_id, event_name, author_id, date, created_at, updated_at
`

type UpdateEventParams struct {
	EventID   uuid.UUID
	Date      time.Time
	EventName string
	UpdatedAt time.Time
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, updateEvent,
		arg.EventID,
		arg.Date,
		arg.EventName,
		arg.UpdatedAt,
	)
	var i Event
	err := row.Scan(
		&i.EventID,
		&i.EventName,
		&i.AuthorID,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
