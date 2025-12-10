package routes

import (
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, authService auth.AuthService) {
	app.Post("/auth/validate", handlers.Validate(authService))
}
