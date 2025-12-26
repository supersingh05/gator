-- name: CreateFeed :one
INSERT INTO feeds(name, url, created_by)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1 LIMIT 1;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = $1;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds
WHERE id is not null;

-- name: GetAllFeeds :many
SELECT * FROM feeds;


-- name: GetAllFeedsWithUserName :many
SELECT f.id, f.name as feedName, f.url, u.name as userName 
FROM feeds AS f
INNER JOIN users AS u
    ON f.created_by = u.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1 LIMIT 1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;
