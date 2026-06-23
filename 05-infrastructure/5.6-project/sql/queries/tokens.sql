-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token_hash, is_revoked, revoked_at, expires_at;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET is_revoked = true, revoked_at = NOW() WHERE is_revoked = false AND token_hash = $1;

-- name: FindActiveRefreshToken :one
SELECT * FROM refresh_tokens WHERE is_revoked = false AND token_hash = $1 AND expires_at > NOW() LIMIT 1;