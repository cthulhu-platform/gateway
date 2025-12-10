-- File repository queries for sqlc
-- Note: Currently using manual SQL in repository implementation
-- These queries can be added if you want to use sqlc-generated code

-- Example queries (commented out until needed):
-- name: GetBucketByID :one
-- SELECT * FROM buckets WHERE id = ? LIMIT 1;

-- name: CreateBucket :exec
-- INSERT INTO buckets (id, password_hash, created_at, updated_at)
-- VALUES (?, ?, ?, ?);

-- name: GetFileByStringID :one
-- SELECT * FROM files WHERE string_id = ? LIMIT 1;

-- name: GetFilesByBucketID :many
-- SELECT * FROM files WHERE bucket_id = ? ORDER BY created_at ASC;

-- name: CreateFile :exec
-- INSERT INTO files (string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at)
-- VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: AddBucketAdmin :exec
-- INSERT INTO bucket_admins (user_id, bucket_id, created_at)
-- VALUES (?, ?, ?);

