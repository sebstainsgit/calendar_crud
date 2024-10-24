-- +goose Up
CREATE TABLE event_users (
    event_id UUID NOT NULL REFERENCES events(event_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, user_id)
);
-- +goose Down
DROP TABLE event_users;