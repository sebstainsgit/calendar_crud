-- name: GetEvents :many
SELECT * FROM events
ORDER BY date DESC;

-- name: GetUsersEvents :many
SELECT * FROM events
WHERE users_id = $1;

-- name: CreateEvent :one
INSERT INTO events (event_id, event_name, users_id, date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE event_id = $1;

-- name: DeleteEventByID :exec
DELETE FROM events
WHERE event_id = $1;