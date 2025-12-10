package handlers

import (
	"fmt"
	"strings"

	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
)

func FileUpload(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid multipart payload",
			})
		}

		files := form.File["files"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "no files provided; expected field 'files'",
			})
		}

		res, err := s.UploadFiles(c.UserContext(), files)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func DownloadFile(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := c.Params("id")
		filename := c.Params("filename")
		if storageID == "" || filename == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "storage id and filename required",
			})
		}

		res, err := s.DownloadFile(c.UserContext(), storageID, filename)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		contentType := res.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Set("Content-Type", contentType)
		if res.ContentLength > 0 {
			c.Set("Content-Length", fmt.Sprintf("%d", res.ContentLength))
		}
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))

		if err := c.SendStream(res.Body); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func RetrieveFileBucket(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := strings.TrimSpace(c.Params("id"))
		if storageID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "storage id required",
			})
		}

		meta, err := s.RetrieveFileBucket(c.UserContext(), storageID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(meta)
	}
}
