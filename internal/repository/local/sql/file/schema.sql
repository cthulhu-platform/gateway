-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Buckets table: Storage containers for files
CREATE TABLE IF NOT EXISTS buckets (
    id TEXT PRIMARY KEY,  -- storage_id/session_id (e.g., "samplebuck", 10-char alphanumeric)
    password_hash TEXT,  -- NULL = public/anonymous access, set = protected (hashing logic deferred)
    created_at INTEGER NOT NULL,  -- Unix timestamp
    updated_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_buckets_created_at ON buckets(created_at);

-- Files table: File metadata and references
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- Numeric primary key
    string_id TEXT NOT NULL UNIQUE,  -- Surrogate key stored in DB, used in S3 path (e.g., "hashid1")
    bucket_id TEXT NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    original_name TEXT NOT NULL,  -- Original filename (e.g., "test.txt")
    owner_id TEXT REFERENCES users(id) ON DELETE SET NULL,  -- Nullable owner reference to auth.users
    size INTEGER NOT NULL,  -- File size in bytes
    content_type TEXT NOT NULL,  -- MIME type
    s3_key TEXT NOT NULL,  -- Full S3 key (e.g., "samplebuck/hashid1")
    created_at INTEGER NOT NULL  -- Unix timestamp
);

CREATE INDEX IF NOT EXISTS idx_files_bucket_id ON files(bucket_id);
CREATE INDEX IF NOT EXISTS idx_files_string_id ON files(string_id);
CREATE INDEX IF NOT EXISTS idx_files_owner_id ON files(owner_id);
CREATE INDEX IF NOT EXISTS idx_files_s3_key ON files(s3_key);

-- Bucket admins table: Many-to-many relationship between users and buckets
CREATE TABLE IF NOT EXISTS bucket_admins (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bucket_id TEXT NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    created_at INTEGER NOT NULL,  -- Unix timestamp
    PRIMARY KEY (user_id, bucket_id)
);

CREATE INDEX IF NOT EXISTS idx_bucket_admins_user_id ON bucket_admins(user_id);
CREATE INDEX IF NOT EXISTS idx_bucket_admins_bucket_id ON bucket_admins(bucket_id);

