-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1,
        NOW(),
        NOW(),
        $2,
        $3,
        $4) RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE refresh_tokens.token = $1;

-- name: UpdateRefreshTokenRevokedAtByToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(),
    revoked_at = $1
WHERE token = $2;
