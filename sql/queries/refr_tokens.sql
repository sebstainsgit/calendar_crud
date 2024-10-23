-- name: CreateRefrToken :one
INSERT INTO refresh_tokens (refr_token, users_id, expires)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UserIDFromRefrToken :one
SELECT users_id FROM refresh_tokens
WHERE refr_token = $1;