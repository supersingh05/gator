-- name: CreatePost :one
INSERT INTO posts (published_at, title, description, url, created_by, feed_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.*
FROM posts p
JOIN user_feeds uf ON uf.feed_id = p.feed_id
WHERE uf.user_id = $1
ORDER BY p.published_at DESC NULLS LAST;

