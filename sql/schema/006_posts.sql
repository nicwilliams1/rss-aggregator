-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title Text NOT NULL,
  url Text NOT NULL UNIQUE,
  description Text,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID REFERENCES feeds (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;