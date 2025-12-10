-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Users table: Core user identity (bare essentials)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,  -- UUID
    oauth_provider TEXT NOT NULL,  -- 'github', 'google', etc.
    oauth_user_id TEXT NOT NULL,  -- Provider's user ID
    email TEXT UNIQUE NOT NULL,
    username TEXT,
    avatar_url TEXT,
    created_at INTEGER NOT NULL,  -- Unix timestamp
    updated_at INTEGER NOT NULL,
    deleted_at INTEGER,  -- NULL if active, Unix timestamp if soft deleted
    UNIQUE(oauth_provider, oauth_user_id)
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_oauth_lookup ON users(oauth_provider, oauth_user_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Refresh tokens table: Simple refresh token storage
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id TEXT PRIMARY KEY,  -- UUID
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,  -- SHA-256 hash
    expires_at INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    revoked_at INTEGER,  -- NULL if active, Unix timestamp if revoked
    revoked_reason TEXT  -- 'user_logout', 'expired', etc.
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_active ON refresh_tokens(user_id, revoked_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);

-- OAuth sessions table: OAuth flow state management
CREATE TABLE IF NOT EXISTS oauth_sessions (
    state TEXT PRIMARY KEY,  -- UUID
    provider TEXT NOT NULL,
    code_verifier TEXT NOT NULL,  -- PKCE
    code_challenge TEXT NOT NULL,
    redirect_uri TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_oauth_sessions_expires ON oauth_sessions(expires_at);

