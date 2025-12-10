package file

import (
	"fmt"
	"time"

	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/golang-jwt/jwt/v5"
)

const (
	bucketTokenExpiry = 30 * time.Minute
)

// BucketAccessClaims represents JWT claims for bucket access tokens
type BucketAccessClaims struct {
	BucketID    string   `json:"bucket_id"`
	Privileges  []string `json:"privileges"` // ["read", "write", etc.]
	UserID      *string  `json:"user_id,omitempty"`
	AuthTokenID *string  `json:"auth_token_id,omitempty"` // JTI from auth token
	jwt.RegisteredClaims
}

// BucketAccessTokenResponse represents the API response for bucket authentication
type BucketAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

// GenerateBucketAccessToken generates a JWT token for bucket access
func GenerateBucketAccessToken(bucketID string, userID *string, authTokenID *string, privileges []string) (string, error) {
	if bucketID == "" {
		return "", fmt.Errorf("bucket_id is required")
	}

	jwtSecret := pkg.JWT_SECRET
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is required")
	}

	now := time.Now()
	claims := &BucketAccessClaims{
		BucketID:    bucketID,
		Privileges:  privileges,
		UserID:      userID,
		AuthTokenID: authTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(bucketTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign bucket access token: %w", err)
	}

	return tokenString, nil
}

// ValidateBucketAccessToken validates a bucket access token and returns its claims
func ValidateBucketAccessToken(tokenString string) (*BucketAccessClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is required")
	}

	jwtSecret := pkg.JWT_SECRET
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	claims := &BucketAccessClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
