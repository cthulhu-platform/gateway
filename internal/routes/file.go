package routes

import (
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
)

func FileRouter(app fiber.Router, fileService file.FileService) {
	app.Post("/files/upload", handlers.FileUpload(fileService))
	app.Get("/files/s/:id", handlers.RetrieveFileBucket(fileService))
	app.Get("/files/s/:id/d/:filename", handlers.DownloadFile(fileService))
}
