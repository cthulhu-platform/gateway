package routes

import (
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/middleware"
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
)

func FileRouter(app fiber.Router, fileService file.FileService, authService auth.AuthService) {
	// Upload route with optional auth middleware
	app.Post("/files/upload", middleware.OptionalJWTAuth(authService), handlers.FileUpload(fileService))
	app.Get("/files/s/:id", handlers.RetrieveFileBucket(fileService))
	app.Get("/files/s/:id/admins", handlers.GetBucketAdmins(fileService))
	app.Get("/files/s/:id/d/:filename", handlers.DownloadFile(fileService))
}
