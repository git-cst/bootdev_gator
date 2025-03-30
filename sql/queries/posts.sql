-- name: CreatePost :one
INSERT INTO posts(
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT 
    p.title,
    p.url,
    p.description,
    p.published_at
FROM posts as p
INNER JOIN feed_follows as ff
ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;