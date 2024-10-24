-- name: CreateUser :one
INSERT INTO users (user_id, name, email, password, created_at, updated_at, elevation)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserFromID :one
SELECT * FROM users
WHERE user_id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, password = $4, updated_at = $5
WHERE user_id = $1
RETURNING *;