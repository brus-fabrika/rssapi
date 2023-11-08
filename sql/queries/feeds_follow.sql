-- name: CreateFeedFollow :one
INSERT INTO feeds_follow (id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedsFollowByUserId :many
SELECT * FROM feeds_follow WHERE user_id = $1;