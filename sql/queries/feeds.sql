-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsWithUserName :many
SELECT feeds.id,
    feeds.created_at,
    feeds.updated_at,
    feeds.name,
    url,
    user_id,
    users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;