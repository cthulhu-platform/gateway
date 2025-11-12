package handlers

import (
	"github.com/cthulhu-platform/gateway/internal/service"
	"github.com/gofiber/fiber/v2"
)

func FileUpload(services *service.ServiceContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fm := services.FileManager

		return c.JSON(fiber.Map{
			"msg": "Uploading file",
		})
	}
}

func DownloadFile(services *service.ServiceContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg": "Downloading file",
		})
	}
}

func RetrieveFileBucket(services *service.ServiceContainer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg": "Retrieving File Bucket",
		})
	}
}
