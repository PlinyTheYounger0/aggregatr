-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT *
FROM feeds;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
        id, 
        created_at, 
        updated_at, 
        user_id, 
        feed_id
    )
    VALUES (
        $1, 
        $2, 
        $3, 
        $4, 
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users
    ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds
    ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetFeedNameByID :one
SELECT name
FROM feeds
WHERE id = $1;

-- name: GetFeedFollowsByUser :many
SELECT *
FROM feed_follows
WHERE user_id = $1;
    