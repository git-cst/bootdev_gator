-- name: CreateFeed :one
INSERT INTO feed(created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeeds :many
SELECT 
    f.name as Name,
    f.url as Url,
    u.name as User,
    f.created_at as Created_at,
    f.updated_at as updated_at
FROM feed as f
INNER JOIN users as u
ON u.id = f.user_id; 

-- name: GetFeedByUrl :one
SELECT 
    f.id as ID,
    f.name as Name,
    f.url as Url,
    u.name as User,
    f.created_at as Created_at,
    f.updated_at as updated_at
FROM feed as f
INNER JOIN users as u
ON u.id = f.user_id
WHERE f.url = $1
LIMIT 1; 

-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows(created_at, updated_at, user_id, feed_id)
    VALUES (
        $1, $2, $3, $4
    )
    RETURNING *
)
SELECT
    inserted_feed_follows.*,
    feed.name AS feed_name,
    users.name AS username
FROM inserted_feed_follows
INNER JOIN feed
ON feed.id = inserted_feed_follows.feed_id
INNER JOIN users
ON users.id = inserted_feed_follows.user_id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    f.name as feed_name,
    u.name as username
FROM feed_follows
INNER JOIN feed as f
ON f.id = feed_follows.feed_id
INNER JOIN users as u
ON u.id = feed_follows.user_id
WHERE feed_follows.user_id = $1;

-- name: RemoveFollowForUser :exec
DELETE FROM
feed_follows
WHERE 
feed_follows.user_id = $1 AND
feed_follows.feed_id = $2;