-- +goose Up
CREATE TABLE refresh_tokens (
    refr_token VARCHAR(64) NOT NULL PRIMARY KEY,
    users_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    expires TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE refresh_tokens;