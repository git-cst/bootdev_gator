-- name: CreateUser :one
INSERT INTO users(created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1
LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: ResetUsers :exec
TRUNCATE TABLE users;