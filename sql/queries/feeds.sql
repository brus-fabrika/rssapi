-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsByUserId :many
SELECT * FROM feeds
WHERE user_id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC
LIMIT 1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;