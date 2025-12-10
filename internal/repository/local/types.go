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
