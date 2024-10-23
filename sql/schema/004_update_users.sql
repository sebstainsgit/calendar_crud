-- +goose Up
ALTER TABLE users
ADD elevation text
NOT NULL DEFAULT 'user';

-- +goose Down
ALTER TABLE users
DROP elevation;