-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2,$3, $4, $5)
RETURNING *;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = $2, revoked_at = $2
WHERE token = $1
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

