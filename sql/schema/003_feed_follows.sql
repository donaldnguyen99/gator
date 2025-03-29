-- +goose Up
CREATE TABLE feed_follows(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT feed_follows_user_id_feed_id_key UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;