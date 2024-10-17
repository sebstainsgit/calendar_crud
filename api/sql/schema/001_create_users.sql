-- +goose Up
CREATE TABLE users (
    user_id UUID NOT NULL PRIMARY KEY,
    api_key VARCHAR(64) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;