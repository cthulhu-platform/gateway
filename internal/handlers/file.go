package handlers

import (
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
)

func FileUpload(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fm := services.FileManager
		s.UploadFile()

		return c.JSON(fiber.Map{
			"msg": "Uploading file",
		})
	}
}

func DownloadFile(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg": "Downloading file",
		})
	}
}

func RetrieveFileBucket(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg": "Retrieving File Bucket",
		})
	}
}
