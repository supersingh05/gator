
-- name: CreateUserFeed :one
INSERT INTO user_feeds(user_id, feed_id)
VALUES($1, $2)
RETURNING *;

-- name: GetFeedsForUser :many
SELECT f.*
FROM feeds AS f
INNER JOIN user_feeds AS uf
    ON f.id = uf.feed_id
WHERE uf.user_id = $1;

-- name: DeleteWholeFeedForUser :exec
DELETE from user_feeds
WHERE user_id = $1;

-- name: DeleteFeedForUser :exec
DELETE FROM user_feeds
WHERE user_id = $1 and feed_id = $2;

