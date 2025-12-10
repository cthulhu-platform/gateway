package handlers

import (
	"fmt"
	"strings"

	"github.com/cthulhu-platform/gateway/internal/microservices/authentication"
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

		// Extract password from form data (optional)
		var password *string
		if passwordValues := form.Value["password"]; len(passwordValues) > 0 && passwordValues[0] != "" {
			password = &passwordValues[0]
		}

		// Extract user_id from context (optional, may be nil)
		var userID *string
		if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
			userID = &uid
		}

		res, err := s.UploadFiles(c.UserContext(), files, userID, password)
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
		stringID := c.Params("filename") // URL param is actually string_id, not filename
		if storageID == "" || stringID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "storage id and string id required",
			})
		}

		res, err := s.DownloadFile(c.UserContext(), storageID, stringID)
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
		// Use the original filename from the download result
		filename := res.DownloadedFile
		if filename == "" {
			filename = stringID
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

func GetBucketAdmins(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := strings.TrimSpace(c.Params("id"))
		if storageID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "storage id required",
			})
		}

		admins, err := s.GetBucketAdmins(c.UserContext(), storageID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(admins)
	}
}

func IsProtected(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := strings.TrimSpace(c.Params("id"))
		if storageID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "storage id required",
			})
		}

		isProtected, _, err := s.IsBucketProtected(c.UserContext(), storageID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"protected": isProtected,
			"bucket_id": storageID,
		})
	}
}

func AuthenticateBucket(s file.FileService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storageID := strings.TrimSpace(c.Params("id"))
		if storageID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "storage id required",
			})
		}

		// Extract password from request body
		var body struct {
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid request body",
			})
		}

		if body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "password is required",
			})
		}

		// Extract user_id from context (if authenticated)
		var userID *string
		if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
			userID = &uid
		}

		// Extract auth token JTI from claims (if available)
		// JTI is in RegisteredClaims.ID field
		var authTokenID *string
		if claims, ok := c.Locals("claims").(*authentication.Claims); ok && claims != nil {
			if claims.ID != "" {
				authTokenID = &claims.ID
			}
		}

		// Authenticate and get token
		token, err := s.AuthenticateBucket(c.UserContext(), storageID, body.Password, userID, authTokenID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"access_token": token,
			"expires_in":   1800, // 30 minutes in seconds
		})
	}
}
