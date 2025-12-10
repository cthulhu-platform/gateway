package routes

import (
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, authService auth.AuthService) {
	// OAuth endpoints
	app.Get("/auth/oauth/:provider", handlers.OAuthInitiate(authService))
	app.Get("/auth/oauth/:provider/callback", handlers.OAuthCallback(authService))
	
	// Token management
	app.Post("/auth/refresh", handlers.RefreshToken(authService))
	app.Post("/auth/logout", handlers.Logout(authService))
	app.Post("/auth/validate", handlers.Validate(authService))
}
