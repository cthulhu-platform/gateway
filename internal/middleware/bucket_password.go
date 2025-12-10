package middleware

import (
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
)

// Note: VerifyPassword is exported from the file package (password.go)
// so we can call it directly

// BucketPasswordAuth middleware verifies bucket password if bucket is protected
func BucketPasswordAuth(fileService file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := c.Params("id")
		if storageID == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "storage id required",
			})
		}

		// Get bucket to check if it's protected
		isProtected, passwordHash, err := fileService.IsBucketProtected(c.UserContext(), storageID)
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

		// If protected, verify password
		providedPassword := c.Get("X-Bucket-Password")
		if providedPassword == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "password required for protected bucket",
			})
		}

		// Verify password using argon2
		if !verifyPassword(providedPassword, *passwordHash) {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "invalid password",
			})
		}

		return c.Next()
	}
}

// verifyPassword verifies a password against a stored hash
// This calls the exported VerifyPassword function from the file package
func verifyPassword(password, hash string) bool {
	return file.VerifyPassword(password, hash)
}
