package local

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	"github.com/cthulhu-platform/gateway/internal/pkg"
	_ "modernc.org/sqlite"
)

//go:embed sql/file/schema.sql
var fileSchemaFS embed.FS

type FileRepository interface {
	GetDB() *sql.DB
	Close() error
	// Bucket operations
	CreateBucket(bucket *Bucket) error
	GetBucketByID(bucketID string) (*Bucket, error)
	UpdateBucket(bucket *Bucket) error
	// File operations
	CreateFile(file *File) error
	GetFileByID(id int64) (*File, error)
	GetFileByStringID(stringID string) (*File, error)
	GetFilesByBucketID(bucketID string) ([]*File, error)
	GetFileByBucketIDAndOriginalName(bucketID, originalName string) (*File, error)
	// Bucket admin operations
	AddBucketAdmin(admin *BucketAdmin) error
	RemoveBucketAdmin(userID, bucketID string) error
	GetBucketAdminsByBucketID(bucketID string) ([]*BucketAdmin, error)
	GetBucketAdminsByUserID(userID string) ([]*BucketAdmin, error)
	IsBucketAdmin(userID, bucketID string) (bool, error)
}

type localFileRepository struct {
	db *sql.DB
}

func NewLocalFileRepository() (*localFileRepository, error) {
	path := pkg.LOCAL_FILE_REPO

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open SQLite database connection
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Initialize schema
	schema, err := fileSchemaFS.ReadFile("sql/file/schema.sql")
	if err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, err
	}

	r := &localFileRepository{
		db: db,
	}

	return r, nil
}

func (r *localFileRepository) GetDB() *sql.DB {
	return r.db
}

func (r *localFileRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// Bucket operations

func (r *localFileRepository) CreateBucket(bucket *Bucket) error {
	query := `INSERT INTO buckets (id, password_hash, created_at, updated_at)
	          VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, bucket.ID, bucket.PasswordHash, bucket.CreatedAt, bucket.UpdatedAt)
	return err
}

func (r *localFileRepository) GetBucketByID(bucketID string) (*Bucket, error) {
	query := `SELECT id, password_hash, created_at, updated_at
	          FROM buckets WHERE id = ? LIMIT 1`

	bucket := &Bucket{}
	var passwordHash sql.NullString

	err := r.db.QueryRow(query, bucketID).Scan(
		&bucket.ID, &passwordHash, &bucket.CreatedAt, &bucket.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if passwordHash.Valid {
		bucket.PasswordHash = &passwordHash.String
	}

	return bucket, nil
}

func (r *localFileRepository) UpdateBucket(bucket *Bucket) error {
	query := `UPDATE buckets SET password_hash = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.Exec(query, bucket.PasswordHash, bucket.UpdatedAt, bucket.ID)
	return err
}

// File operations

func (r *localFileRepository) CreateFile(file *File) error {
	query := `INSERT INTO files (string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		file.StringID, file.BucketID, file.OriginalName, file.OwnerID,
		file.Size, file.ContentType, file.S3Key, file.CreatedAt,
	)
	return err
}

func (r *localFileRepository) GetFileByID(id int64) (*File, error) {
	query := `SELECT id, string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at
	          FROM files WHERE id = ? LIMIT 1`

	file := &File{}
	var ownerID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&file.ID, &file.StringID, &file.BucketID, &file.OriginalName,
		&ownerID, &file.Size, &file.ContentType, &file.S3Key, &file.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if ownerID.Valid {
		file.OwnerID = &ownerID.String
	}

	return file, nil
}

func (r *localFileRepository) GetFileByStringID(stringID string) (*File, error) {
	query := `SELECT id, string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at
	          FROM files WHERE string_id = ? LIMIT 1`

	file := &File{}
	var ownerID sql.NullString

	err := r.db.QueryRow(query, stringID).Scan(
		&file.ID, &file.StringID, &file.BucketID, &file.OriginalName,
		&ownerID, &file.Size, &file.ContentType, &file.S3Key, &file.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if ownerID.Valid {
		file.OwnerID = &ownerID.String
	}

	return file, nil
}

func (r *localFileRepository) GetFilesByBucketID(bucketID string) ([]*File, error) {
	query := `SELECT id, string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at
	          FROM files WHERE bucket_id = ? ORDER BY created_at ASC`

	rows, err := r.db.Query(query, bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*File, 0)
	for rows.Next() {
		file := &File{}
		var ownerID sql.NullString

		err := rows.Scan(
			&file.ID, &file.StringID, &file.BucketID, &file.OriginalName,
			&ownerID, &file.Size, &file.ContentType, &file.S3Key, &file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if ownerID.Valid {
			file.OwnerID = &ownerID.String
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *localFileRepository) GetFileByBucketIDAndOriginalName(bucketID, originalName string) (*File, error) {
	query := `SELECT id, string_id, bucket_id, original_name, owner_id, size, content_type, s3_key, created_at
	          FROM files WHERE bucket_id = ? AND original_name = ? LIMIT 1`

	file := &File{}
	var ownerID sql.NullString

	err := r.db.QueryRow(query, bucketID, originalName).Scan(
		&file.ID, &file.StringID, &file.BucketID, &file.OriginalName,
		&ownerID, &file.Size, &file.ContentType, &file.S3Key, &file.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if ownerID.Valid {
		file.OwnerID = &ownerID.String
	}

	return file, nil
}

// Bucket admin operations

func (r *localFileRepository) AddBucketAdmin(admin *BucketAdmin) error {
	query := `INSERT INTO bucket_admins (user_id, bucket_id, created_at)
	          VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, admin.UserID, admin.BucketID, admin.CreatedAt)
	return err
}

func (r *localFileRepository) RemoveBucketAdmin(userID, bucketID string) error {
	query := `DELETE FROM bucket_admins WHERE user_id = ? AND bucket_id = ?`

	_, err := r.db.Exec(query, userID, bucketID)
	return err
}

func (r *localFileRepository) GetBucketAdminsByBucketID(bucketID string) ([]*BucketAdmin, error) {
	query := `SELECT user_id, bucket_id, created_at
	          FROM bucket_admins WHERE bucket_id = ? ORDER BY created_at ASC`

	rows, err := r.db.Query(query, bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	admins := make([]*BucketAdmin, 0)
	for rows.Next() {
		admin := &BucketAdmin{}
		err := rows.Scan(&admin.UserID, &admin.BucketID, &admin.CreatedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
}

func (r *localFileRepository) GetBucketAdminsByUserID(userID string) ([]*BucketAdmin, error) {
	query := `SELECT user_id, bucket_id, created_at
	          FROM bucket_admins WHERE user_id = ? ORDER BY created_at ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	admins := make([]*BucketAdmin, 0)
	for rows.Next() {
		admin := &BucketAdmin{}
		err := rows.Scan(&admin.UserID, &admin.BucketID, &admin.CreatedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
}

func (r *localFileRepository) IsBucketAdmin(userID, bucketID string) (bool, error) {
	query := `SELECT 1 FROM bucket_admins WHERE user_id = ? AND bucket_id = ? LIMIT 1`

	var exists int
	err := r.db.QueryRow(query, userID, bucketID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
