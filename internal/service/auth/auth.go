package auth

import "github.com/cthulhu-platform/gateway/internal/microservices/authentication"

type AuthService interface {
	InitiateOAuth(provider string) (string, error)
	HandleOAuthCallback(provider, code, state string) (*AuthResponse, error)
	ValidateToken(token string) (*authentication.Claims, error)
	RefreshToken(refreshToken string) (*authentication.TokenPair, error)
	Logout(refreshToken string) error
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Username *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
