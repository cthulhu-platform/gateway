package handlers

import (
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

// OAuthInitiate initiates OAuth flow for a provider
func OAuthInitiate(s auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Params("provider")
		if provider == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "provider parameter is required",
			})
		}

		oauthURL, err := s.InitiateOAuth(provider)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Redirect(oauthURL, 302)
	}
}

// OAuthCallback handles OAuth callback
func OAuthCallback(s auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Params("provider")
		code := c.Query("code")
		state := c.Query("state")

		if provider == "" || code == "" || state == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "provider, code, and state are required",
			})
		}

		authResponse, err := s.HandleOAuthCallback(provider, code, state)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(authResponse)
	}
}

// RefreshToken handles token refresh
func RefreshToken(s auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if req.RefreshToken == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "refresh_token is required",
			})
		}

		tokenPair, err := s.RefreshToken(req.RefreshToken)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(tokenPair)
	}
}

// Logout handles user logout
func Logout(s auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if req.RefreshToken == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "refresh_token is required",
			})
		}

		if err := s.Logout(req.RefreshToken); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "logged out successfully",
		})
	}
}

// Validate validates an access token
func Validate(s auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "authorization header is required",
			})
		}

		// Extract token from "Bearer <token>"
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		claims, err := s.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"valid": true,
			"claims": claims,
		})
	}
}
