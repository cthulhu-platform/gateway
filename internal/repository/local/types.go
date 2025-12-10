package local

// User represents a user in the system
type User struct {
	ID            string
	OAuthProvider string
	OAuthUserID   string
	Email         string
	Username      *string
	AvatarURL     *string
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     *int64
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	ID            string
	UserID        string
	TokenHash     string
	ExpiresAt     int64
	CreatedAt     int64
	RevokedAt     *int64
	RevokedReason *string
}

// OAuthSession represents an OAuth flow session
type OAuthSession struct {
	State         string
	Provider      string
	CodeVerifier  string
	CodeChallenge string
	RedirectURI   string
	ExpiresAt     int64
	CreatedAt     int64
}

// Bucket represents a storage container for files
type Bucket struct {
	ID           string  // storage_id/session_id (e.g., "samplebuck")
	PasswordHash *string // NULL = public/anonymous, set = protected
	CreatedAt    int64
	UpdatedAt    int64
}

// File represents file metadata
type File struct {
	ID           int64   // Numeric primary key
	StringID     string  // Surrogate key used in S3 path (e.g., "hashid1")
	BucketID     string  // References buckets(id)
	OriginalName string  // Original filename (e.g., "test.txt")
	OwnerID      *string // Nullable owner reference to users(id)
	Size         int64   // File size in bytes
	ContentType  string  // MIME type
	S3Key        string  // Full S3 key (e.g., "samplebuck/hashid1")
	CreatedAt    int64
}

// BucketAdmin represents a many-to-many relationship between users and buckets
type BucketAdmin struct {
	UserID    string
	BucketID  string
	CreatedAt int64
}
