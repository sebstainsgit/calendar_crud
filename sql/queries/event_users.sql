-- name: AddUserToEvent :one
INSERT INTO event_users (event_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserIDsForEvent :many
SELECT u.user_id
FROM event_users eu
JOIN users u ON eu.user_id = u.user_id
WHERE eu.event_id = $1;

-- name: GetEventsIDsForUser :many
SELECT event_id
FROM event_users 
WHERE user_id = $1;

-- name: DeleteUserFromEvent :exec
DELETE FROM event_users
WHERE event_id = $1 AND user_id = $2;

-- name: GetEventConcerns :one
SELECT * FROM event_users
WHERE event_id = $1 AND user_id = $2;