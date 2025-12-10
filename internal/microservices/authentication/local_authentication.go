package authentication

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/repository/local"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LocalAuthenticationConnection struct {
	repo          local.AuthRepository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewLocalAuthConnection() (*LocalAuthenticationConnection, error) {
	repo, err := local.NewLocalAuthRepository()
	if err != nil {
		return nil, err
	}

	jwtSecret := pkg.JWT_SECRET
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	accessExpiry, err := time.ParseDuration(pkg.JWT_ACCESS_EXPIRY)
	if err != nil {
		accessExpiry = 15 * time.Minute // Default
	}

	refreshExpiry, err := time.ParseDuration(pkg.JWT_REFRESH_EXPIRY)
	if err != nil {
		refreshExpiry = 7 * 24 * time.Hour // Default 7 days
	}

	c := &LocalAuthenticationConnection{
		repo:          repo,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
	return c, nil
}

func (c *LocalAuthenticationConnection) Close() {
	if c.repo != nil {
		c.repo.Close()
	}
}

// GetRepo returns the repository for direct access (used by service layer)
func (c *LocalAuthenticationConnection) GetRepo() local.AuthRepository {
	return c.repo
}

func (c *LocalAuthenticationConnection) GenerateTokens(userID, email, provider string) (*TokenPair, error) {
	now := time.Now()

	// Generate access token
	accessClaims := &Claims{
		UserID:   userID,
		Email:    email,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(c.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(c.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token (random string, not JWT)
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	// Hash refresh token for storage
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Store refresh token in database
	refreshTokenID := uuid.New().String()
	refreshTokenRecord := &local.RefreshToken{
		ID:        refreshTokenID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(c.refreshExpiry).Unix(),
		CreatedAt: now.Unix(),
	}

	if err := c.repo.CreateRefreshToken(refreshTokenRecord); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
	}, nil
}

func (c *LocalAuthenticationConnection) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (c *LocalAuthenticationConnection) ValidateUserID(userID string) (bool, error) {
	if userID == "" {
		return false, nil
	}
	user, err := c.repo.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
