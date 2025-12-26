-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * from users
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * from users
WHERE name = $1 LIMIT 1;

-- name: DeleteAllUsers :exec
DELETE FROM users where id IS NOT NULL;

-- name: GetUsers :many
SELECT * from users;
