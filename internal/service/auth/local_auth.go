package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/cthulhu-platform/gateway/internal/microservices/authentication"
	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/repository/local"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type localAuthService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewLocalAuthService(conns *microservices.ServiceConnectionContainer) *localAuthService {
	return &localAuthService{
		conns: conns,
	}
}

func (s *localAuthService) InitiateOAuth(provider string) (string, error) {
	// Get authentication connection
	authConn, ok := s.conns.Authentication.(*authentication.LocalAuthenticationConnection)
	if !ok {
		return "", fmt.Errorf("authentication connection is not local")
	}

	// Get repository from connection
	repo := authConn.GetRepo()
	if repo == nil {
		return "", fmt.Errorf("repository not available")
	}

	// Generate PKCE values
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Generate state
	state := uuid.New().String()

	// Store OAuth session
	session := &local.OAuthSession{
		State:         state,
		Provider:      provider,
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
		RedirectURI:   getRedirectURI(provider),
		ExpiresAt:     time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:     time.Now().Unix(),
	}

	if err := repo.CreateOAuthSession(session); err != nil {
		return "", fmt.Errorf("failed to create OAuth session: %w", err)
	}

	// Build OAuth URL
	oauthURL := buildOAuthURL(provider, state, codeChallenge)
	return oauthURL, nil
}

func (s *localAuthService) HandleOAuthCallback(provider, code, state string) (*AuthResponse, error) {
	// Get authentication connection
	authConn, ok := s.conns.Authentication.(*authentication.LocalAuthenticationConnection)
	if !ok {
		return nil, fmt.Errorf("authentication connection is not local")
	}

	repo := authConn.GetRepo()
	if repo == nil {
		return nil, fmt.Errorf("repository not available")
	}

	// Validate state
	session, err := repo.GetOAuthSession(state)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("invalid or expired OAuth session")
	}

	if time.Now().Unix() > session.ExpiresAt {
		repo.DeleteOAuthSession(state)
		return nil, fmt.Errorf("OAuth session expired")
	}

	if session.Provider != provider {
		return nil, fmt.Errorf("provider mismatch")
	}

	// Exchange code for token
	oauthConfig := getOAuthConfig(provider)
	token, err := oauthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("code_verifier", session.CodeVerifier))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from provider
	userInfo, err := fetchUserInfo(provider, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Check if user exists
	existingUser, err := repo.GetUserByOAuthID(provider, userInfo.OAuthUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	var user *local.User
	now := time.Now().Unix()

	if existingUser != nil {
		// Update user info
		existingUser.Username = userInfo.Username
		existingUser.AvatarURL = userInfo.AvatarURL
		existingUser.UpdatedAt = now
		if err := repo.UpdateUser(existingUser); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		user = existingUser
	} else {
		// Create new user
		user = &local.User{
			ID:           uuid.New().String(),
			OAuthProvider: provider,
			OAuthUserID:   userInfo.OAuthUserID,
			Email:         userInfo.Email,
			Username:      userInfo.Username,
			AvatarURL:     userInfo.AvatarURL,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := repo.CreateUser(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate tokens
	tokenPair, err := authConn.GenerateTokens(user.ID, user.Email, user.OAuthProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Delete OAuth session (one-time use)
	repo.DeleteOAuthSession(state)

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: &UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			AvatarURL: user.AvatarURL,
		},
	}, nil
}

func (s *localAuthService) ValidateToken(token string) (*authentication.Claims, error) {
	authConn, ok := s.conns.Authentication.(*authentication.LocalAuthenticationConnection)
	if !ok {
		return nil, fmt.Errorf("authentication connection is not local")
	}

	return authConn.ValidateAccessToken(token)
}

func (s *localAuthService) RefreshToken(refreshToken string) (*authentication.TokenPair, error) {
	authConn, ok := s.conns.Authentication.(*authentication.LocalAuthenticationConnection)
	if !ok {
		return nil, fmt.Errorf("authentication connection is not local")
	}

	repo := authConn.GetRepo()
	if repo == nil {
		return nil, fmt.Errorf("repository not available")
	}

	// Hash the refresh token
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Get refresh token from database
	tokenRecord, err := repo.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	if tokenRecord == nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if token is revoked
	if tokenRecord.RevokedAt != nil {
		return nil, fmt.Errorf("refresh token has been revoked")
	}

	// Check if token is expired
	if time.Now().Unix() > tokenRecord.ExpiresAt {
		// Revoke expired token
		repo.RevokeRefreshToken(tokenRecord.ID, "expired")
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := repo.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Revoke old token
	if err := repo.RevokeRefreshToken(tokenRecord.ID, "token_refreshed"); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new tokens
	return authConn.GenerateTokens(user.ID, user.Email, user.OAuthProvider)
}

func (s *localAuthService) Logout(refreshToken string) error {
	authConn, ok := s.conns.Authentication.(*authentication.LocalAuthenticationConnection)
	if !ok {
		return fmt.Errorf("authentication connection is not local")
	}

	repo := authConn.GetRepo()
	if repo == nil {
		return fmt.Errorf("repository not available")
	}

	// Hash the refresh token
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Get refresh token from database
	tokenRecord, err := repo.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		return fmt.Errorf("failed to get refresh token: %w", err)
	}
	if tokenRecord == nil {
		return nil // Token doesn't exist, consider it logged out
	}

	// Revoke token
	return repo.RevokeRefreshToken(tokenRecord.ID, "user_logout")
}

// Helper functions

func generateCodeVerifier() string {
	b := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func buildOAuthURL(provider, state, codeChallenge string) string {
	baseURL := "https://github.com/login/oauth/authorize"
	clientID := pkg.GITHUB_CLIENT_ID
	redirectURI := getRedirectURI(provider)

	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=S256&scope=read:user user:email",
		baseURL, clientID, redirectURI, state, codeChallenge)
}

func getRedirectURI(provider string) string {
	if provider == "github" {
		return pkg.GITHUB_REDIRECT_URI
	}
	return ""
}

func getOAuthConfig(provider string) *oauth2.Config {
	if provider == "github" {
		return &oauth2.Config{
			ClientID:     pkg.GITHUB_CLIENT_ID,
			ClientSecret: pkg.GITHUB_CLIENT_SECRET,
			RedirectURL:  pkg.GITHUB_REDIRECT_URI,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
	}
	return nil
}

type githubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

func fetchUserInfo(provider, accessToken string) (*struct {
	OAuthUserID string
	Email       string
	Username    *string
	AvatarURL   *string
}, error) {
	if provider == "github" {
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch user info: status %d", resp.StatusCode)
		}

		var githubUser githubUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			return nil, err
		}

		// Get email if not in profile
		if githubUser.Email == "" {
			email, err := fetchGitHubEmail(accessToken)
			if err == nil && email != "" {
				githubUser.Email = email
			}
		}

		username := githubUser.Login
		avatarURL := githubUser.AvatarURL

		return &struct {
			OAuthUserID string
			Email       string
			Username    *string
			AvatarURL   *string
		}{
			OAuthUserID: fmt.Sprintf("%d", githubUser.ID),
			Email:       githubUser.Email,
			Username:    &username,
			AvatarURL:   &avatarURL,
		}, nil
	}

	return nil, fmt.Errorf("unsupported provider: %s", provider)
}

func fetchGitHubEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch emails: status %d", resp.StatusCode)
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", nil
}
