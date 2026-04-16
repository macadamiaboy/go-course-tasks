-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token_hash, is_revoked, expires_at;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET is_revoked = true WHERE is_revoked = false AND user_id = $1;

-- name: FindActiveRefreshToken :one
SELECT (id, user_id, token_hash, is_revoked, expires_at) FROM refresh_tokens WHERE is_revoked = false AND user_id = $1 LIMIT 1;