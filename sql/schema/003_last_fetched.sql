-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;

ALTER TABLE feeds
ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();

ALTER TABLE feeds
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

CREATE INDEX idx_feeds_last_fetched_at ON feeds (last_fetched_at);
-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;
ALTER TABLE feeds
DROP COLUMN created_at;
ALTER TABLE feeds
DROP COLUMN updated_at;

DROP INDEX IF EXISTS  idx_feeds_last_fetched_at;
