-- name: GetUserByOAuthID :one
SELECT * FROM users
WHERE oauth_provider = ? AND oauth_user_id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
    id, oauth_provider, oauth_user_id, email, username, avatar_url,
    created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users
SET username = ?, avatar_url = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = ?, updated_at = ?
WHERE id = ?;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    id, user_id, token_hash, expires_at, created_at
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = ? AND revoked_at IS NULL
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = ?, revoked_reason = ?
WHERE id = ?;

-- name: RevokeAllUserTokens :exec
UPDATE refresh_tokens
SET revoked_at = ?, revoked_reason = ?
WHERE user_id = ? AND revoked_at IS NULL;

-- name: CreateOAuthSession :exec
INSERT INTO oauth_sessions (
    state, provider, code_verifier, code_challenge, redirect_uri,
    expires_at, created_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
);

-- name: GetOAuthSession :one
SELECT * FROM oauth_sessions
WHERE state = ?
LIMIT 1;

-- name: DeleteOAuthSession :exec
DELETE FROM oauth_sessions
WHERE state = ?;

-- name: CleanupExpiredOAuthSessions :exec
DELETE FROM oauth_sessions
WHERE expires_at < ?;

