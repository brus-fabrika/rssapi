-- +goose Up
CREATE TABLE feeds_follow (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feeds_follow;
