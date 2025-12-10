package middleware

import (
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

// JWTAuth middleware validates JWT tokens and attaches user claims to context
func JWTAuth(authService auth.AuthService) fiber.Handler {
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

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Attach claims to context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("provider", claims.Provider)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// OptionalJWTAuth middleware validates JWT tokens if present but allows anonymous requests
// If a valid token is provided, it attaches user claims to context
func OptionalJWTAuth(authService auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No auth header, continue without setting user context
			return c.Next()
		}

		// Extract token from "Bearer <token>"
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			// Invalid format, but don't fail - just continue without auth
			return c.Next()
		}

		// Validate token - if invalid, just continue without setting user context
		claims, err := authService.ValidateToken(token)
		if err != nil {
			// Invalid token, but don't fail - just continue without auth
			return c.Next()
		}

		// Valid token - attach claims to context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("provider", claims.Provider)
		c.Locals("claims", claims)

		return c.Next()
	}
}
