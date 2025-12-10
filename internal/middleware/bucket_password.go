package middleware

import (
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
)

// BucketPasswordAuth middleware verifies bucket access token if bucket is protected
// X-Bucket-Password is removed for security - only X-Bucket-Access tokens are accepted
func BucketPasswordAuth(fileService file.FileService, authService auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := c.Params("id")
		if storageID == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "storage id required",
			})
		}

		// Get bucket to check if it's protected
		isProtected, _, err := fileService.IsBucketProtected(c.UserContext(), storageID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   "bucket not found",
			})
		}

		// If not protected, allow access
		if !isProtected {
			return c.Next()
		}

		// If protected, require bucket access token
		accessToken := c.Get("X-Bucket-Access")
		if accessToken == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "bucket access token required",
			})
		}

		// Validate bucket access token
		claims, err := file.ValidateBucketAccessToken(accessToken)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "invalid or expired bucket access token",
			})
		}

		// Verify bucket_id matches
		if claims.BucketID != storageID {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "token bucket_id mismatch",
			})
		}

		// Optional: If token has auth_token_id, validate auth token is still valid
		if claims.AuthTokenID != nil && authService != nil {
			authHeader := c.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				authToken := authHeader[7:]
				authClaims, err := authService.ValidateToken(authToken)
				if err != nil {
					return c.Status(401).JSON(fiber.Map{
						"success": false,
						"error":   "linked auth token is invalid",
					})
				}
				// Verify JTI matches
				if authClaims.ID != *claims.AuthTokenID {
					return c.Status(401).JSON(fiber.Map{
						"success": false,
						"error":   "auth token mismatch",
					})
				}
			}
		}

		return c.Next()
	}
}
