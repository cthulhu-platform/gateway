package routes

import (
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/service"
	"github.com/gofiber/fiber/v2"
)

func FileRouter(app fiber.Router, services *service.ServiceContainer) {
	app.Post("/files/upload", handlers.FileUpload(services))
	app.Get("/files/s/:id", handlers.RetrieveFileBucket(services))
	app.Get("/files/s/:id/d/:filename", handlers.RetrieveFileBucket(services))
}
