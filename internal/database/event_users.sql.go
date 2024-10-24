// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: event_users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addUserToEvent = `-- name: AddUserToEvent :one
INSERT INTO event_users (event_id, user_id)
VALUES ($1, $2)
RETURNING event_id, user_id
`

type AddUserToEventParams struct {
	EventID uuid.UUID
	UserID  uuid.UUID
}

func (q *Queries) AddUserToEvent(ctx context.Context, arg AddUserToEventParams) (EventUser, error) {
	row := q.db.QueryRowContext(ctx, addUserToEvent, arg.EventID, arg.UserID)
	var i EventUser
	err := row.Scan(&i.EventID, &i.UserID)
	return i, err
}

const deleteUserFromEvent = `-- name: DeleteUserFromEvent :exec
DELETE FROM event_users
WHERE event_id = $1 AND user_id = $2
`

type DeleteUserFromEventParams struct {
	EventID uuid.UUID
	UserID  uuid.UUID
}

func (q *Queries) DeleteUserFromEvent(ctx context.Context, arg DeleteUserFromEventParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserFromEvent, arg.EventID, arg.UserID)
	return err
}

const getEventsForUser = `-- name: GetEventsForUser :many
SELECT event_id
FROM event_users 
WHERE user_id = $1
`

func (q *Queries) GetEventsForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getEventsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var event_id uuid.UUID
		if err := rows.Scan(&event_id); err != nil {
			return nil, err
		}
		items = append(items, event_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersForEvent = `-- name: GetUsersForEvent :many
SELECT u.user_id
FROM event_users eu
JOIN users u ON eu.user_id = u.user_id
WHERE eu.event_id = $1
`

func (q *Queries) GetUsersForEvent(ctx context.Context, eventID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getUsersForEvent, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var user_id uuid.UUID
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}