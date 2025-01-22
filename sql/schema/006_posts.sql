-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMPTZ NOT NULL,
    url TEXT NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;