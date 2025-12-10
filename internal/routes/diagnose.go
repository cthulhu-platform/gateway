package routes

import (
	"github.com/cthulhu-platform/gateway/internal/handlers"
	"github.com/cthulhu-platform/gateway/internal/service/diagnose"
	"github.com/gofiber/fiber/v2"
)

func DiagnoseRouter(app fiber.Router, diagnoseService diagnose.DiagnoseService) {
	app.Get("/diagnose/services/all", handlers.ServiceFanoutTest(diagnoseService))
}
