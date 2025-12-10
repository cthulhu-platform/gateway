package handlers

import (
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

func Validate(s auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		s.Validate()

		return c.JSON(fiber.Map{
			"msg": "Validating",
		})
	}
}
