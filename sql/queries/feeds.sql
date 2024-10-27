-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsWithUsername :many
SELECT feeds.name, feeds.url, users.name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = (SELECT CURRENT_TIMESTAMP), updated_at = (SELECT CURRENT_TIMESTAMP)
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;