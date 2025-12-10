package local

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/cthulhu-platform/gateway/internal/pkg"
	_ "modernc.org/sqlite"
)

//go:embed sql/auth/schema.sql
var schemaFS embed.FS

type AuthRepository interface {
	GetDB() *sql.DB
	Close() error
	// User operations
	GetUserByOAuthID(provider, oauthUserID string) (*User, error)
	GetUserByID(userID string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	SoftDeleteUser(userID string) error
	// Refresh token operations
	CreateRefreshToken(token *RefreshToken) error
	GetRefreshTokenByHash(tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(tokenID, reason string) error
	RevokeAllUserTokens(userID, reason string) error
	// OAuth session operations
	CreateOAuthSession(session *OAuthSession) error
	GetOAuthSession(state string) (*OAuthSession, error)
	DeleteOAuthSession(state string) error
}

type localAuthRepository struct {
	db *sql.DB
}

func NewLocalAuthRepository() (*localAuthRepository, error) {
	path := pkg.LOCAL_AUTH_REPO

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
	schema, err := schemaFS.ReadFile("sql/auth/schema.sql")
	if err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, err
	}

	r := &localAuthRepository{
		db: db,
	}

	return r, nil
}

func (r *localAuthRepository) GetDB() *sql.DB {
	return r.db
}

func (r *localAuthRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// User operations

func (r *localAuthRepository) GetUserByOAuthID(provider, oauthUserID string) (*User, error) {
	query := `SELECT id, oauth_provider, oauth_user_id, email, username, avatar_url, created_at, updated_at, deleted_at
	          FROM users WHERE oauth_provider = ? AND oauth_user_id = ? AND deleted_at IS NULL LIMIT 1`

	user := &User{}
	var username, avatarURL sql.NullString
	var deletedAt sql.NullInt64

	err := r.db.QueryRow(query, provider, oauthUserID).Scan(
		&user.ID, &user.OAuthProvider, &user.OAuthUserID, &user.Email,
		&username, &avatarURL, &user.CreatedAt, &user.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if username.Valid {
		user.Username = &username.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Int64
	}

	return user, nil
}

func (r *localAuthRepository) GetUserByID(userID string) (*User, error) {
	query := `SELECT id, oauth_provider, oauth_user_id, email, username, avatar_url, created_at, updated_at, deleted_at
	          FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1`

	user := &User{}
	var username, avatarURL sql.NullString
	var deletedAt sql.NullInt64

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID, &user.OAuthProvider, &user.OAuthUserID, &user.Email,
		&username, &avatarURL, &user.CreatedAt, &user.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if username.Valid {
		user.Username = &username.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Int64
	}

	return user, nil
}

func (r *localAuthRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, oauth_provider, oauth_user_id, email, username, avatar_url, created_at, updated_at, deleted_at
	          FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1`

	user := &User{}
	var username, avatarURL sql.NullString
	var deletedAt sql.NullInt64

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.OAuthProvider, &user.OAuthUserID, &user.Email,
		&username, &avatarURL, &user.CreatedAt, &user.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if username.Valid {
		user.Username = &username.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Int64
	}

	return user, nil
}

func (r *localAuthRepository) CreateUser(user *User) error {
	query := `INSERT INTO users (id, oauth_provider, oauth_user_id, email, username, avatar_url, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		user.ID, user.OAuthProvider, user.OAuthUserID, user.Email,
		user.Username, user.AvatarURL, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *localAuthRepository) UpdateUser(user *User) error {
	query := `UPDATE users SET username = ?, avatar_url = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.Exec(query, user.Username, user.AvatarURL, user.UpdatedAt, user.ID)
	return err
}

func (r *localAuthRepository) SoftDeleteUser(userID string) error {
	query := `UPDATE users SET deleted_at = ?, updated_at = ? WHERE id = ?`

	now := time.Now().Unix()
	_, err := r.db.Exec(query, now, now, userID)
	return err
}

// Refresh token operations

func (r *localAuthRepository) CreateRefreshToken(token *RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
	          VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *localAuthRepository) GetRefreshTokenByHash(tokenHash string) (*RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at, revoked_at, revoked_reason
	          FROM refresh_tokens WHERE token_hash = ? AND revoked_at IS NULL LIMIT 1`

	token := &RefreshToken{}
	var revokedAt sql.NullInt64
	var revokedReason sql.NullString

	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt,
		&token.CreatedAt, &revokedAt, &revokedReason,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Int64
	}
	if revokedReason.Valid {
		token.RevokedReason = &revokedReason.String
	}

	return token, nil
}

func (r *localAuthRepository) RevokeRefreshToken(tokenID, reason string) error {
	query := `UPDATE refresh_tokens SET revoked_at = ?, revoked_reason = ? WHERE id = ?`

	now := time.Now().Unix()
	_, err := r.db.Exec(query, now, reason, tokenID)
	return err
}

func (r *localAuthRepository) RevokeAllUserTokens(userID, reason string) error {
	query := `UPDATE refresh_tokens SET revoked_at = ?, revoked_reason = ? WHERE user_id = ? AND revoked_at IS NULL`

	now := time.Now().Unix()
	_, err := r.db.Exec(query, now, reason, userID)
	return err
}

// OAuth session operations

func (r *localAuthRepository) CreateOAuthSession(session *OAuthSession) error {
	query := `INSERT INTO oauth_sessions (state, provider, code_verifier, code_challenge, redirect_uri, expires_at, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		session.State, session.Provider, session.CodeVerifier, session.CodeChallenge,
		session.RedirectURI, session.ExpiresAt, session.CreatedAt,
	)
	return err
}

func (r *localAuthRepository) GetOAuthSession(state string) (*OAuthSession, error) {
	query := `SELECT state, provider, code_verifier, code_challenge, redirect_uri, expires_at, created_at
	          FROM oauth_sessions WHERE state = ? LIMIT 1`

	session := &OAuthSession{}
	err := r.db.QueryRow(query, state).Scan(
		&session.State, &session.Provider, &session.CodeVerifier, &session.CodeChallenge,
		&session.RedirectURI, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return session, nil
}

func (r *localAuthRepository) DeleteOAuthSession(state string) error {
	query := `DELETE FROM oauth_sessions WHERE state = ?`
	_, err := r.db.Exec(query, state)
	return err
}
