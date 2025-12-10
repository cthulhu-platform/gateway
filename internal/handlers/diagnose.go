package handlers

import (
	"github.com/cthulhu-platform/gateway/internal/service/diagnose"
	"github.com/gofiber/fiber/v2"
)

func ServiceFanoutTest(s diagnose.DiagnoseService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		transactionID, err := s.ServiceFanoutTest()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"transaction_id": transactionID,
			"message":        "Diagnostic message sent to all services",
		})
	}
}
