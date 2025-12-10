package server

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/routes"
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/cthulhu-platform/gateway/internal/service/diagnose"
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
)

type FiberServerConfig struct {
	Host string
	Port string
}

type FiberServer struct {
	Config          *FiberServerConfig
	fileService     file.FileService
	authService     auth.AuthService
	diagnoseService diagnose.DiagnoseService
}

// NOTE: Inject dependencies here
func NewFiberServer(
	c *FiberServerConfig,
	fileService file.FileService,
	authService auth.AuthService,
	diagnoseService diagnose.DiagnoseService,
) *FiberServer {
	return &FiberServer{
		Config:          c,
		fileService:     fileService,
		authService:     authService,
		diagnoseService: diagnoseService,
	}
}

func (s *FiberServer) Start() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// SETUP DEPENDENCIES
	app := fiber.New(fiber.Config{
		BodyLimit: pkg.BODY_LIMIT_MB * 1024 * 1024,
	})

	slog.Info("secret", "github", pkg.GITHUB_CLIENT_ID)

	// SETUP MIDDLEWARE
	app.Use(cors.New(cors.Config{
		AllowOrigins: pkg.CORS_ORIGIN,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(slogfiber.New(logger))
	app.Use(recover.New())

	// ROUTES
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world")
	})
	routes.FileRouter(app, s.fileService)
	routes.AuthRouter(app, s.authService)
	routes.DiagnoseRouter(app, s.diagnoseService)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		app.Shutdown()
	}()

	if err := app.Listen(s.Config.Host + ":" + s.Config.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
