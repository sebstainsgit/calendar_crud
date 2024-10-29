-- name: GetEvents :many
SELECT * FROM events
ORDER BY date DESC;

-- name: GetUserMadeEvents :many
SELECT * FROM events
WHERE author_id = $1;

-- name: CreateEvent :one
INSERT INTO events (event_id, event_name, author_id, date, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE event_id = $1;

-- name: GetEventByName :one
SELECT * FROM events
WHERE event_name = $1;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE event_id = $1;

-- name: UpdateEvent :one
UPDATE events
SET date = $2, event_name = $3, updated_at = $4
WHERE event_id = $1
RETURNING *;