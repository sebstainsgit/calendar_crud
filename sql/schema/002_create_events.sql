-- +goose Up
CREATE TABLE events (
    event_id UUID NOT NULL PRIMARY KEY,
    event_name TEXT NOT NULL,
    users_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE events;