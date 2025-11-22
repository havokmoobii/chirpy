-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days'
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE token = $1;

-- name: GetRefreshTokenRevoked :one
SELECT revoked_at IS NULL
FROM refresh_tokens
WHERE token = $1;

-- name: GetRefreshTokenExpired :one
SELECT NOW() > expires_at
FROM refresh_tokens
WHERE token = $1;