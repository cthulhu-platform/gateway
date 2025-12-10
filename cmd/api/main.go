package main

import (
	"context"
	"fmt"

	"github.com/cthulhu-platform/common/pkg/env"
	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/server"
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/cthulhu-platform/gateway/internal/service/diagnose"
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

func main() {
	// Initialize Viper for environment variable management
	// Supports .env files and environment variables (including APP_ prefix)
	if err := env.Init(".env", "./.env", "../.env"); err != nil {
		// Non-fatal: will fall back to environment variables
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Depency Initialization (RabbitMQ conn, DB conn, SLogger)
	ctx := context.Background()

	// Create RabbitMQ connection with labeled connection name
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		pkg.AMPQ_USER,
		pkg.AMPQ_PASS,
		pkg.AMPQ_HOST,
		pkg.AMPQ_PORT,
		pkg.AMPQ_VHOST,
	)
	conn, err := rabbitmq.NewConn(
		connectionString,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsConfig(rabbitmq.Config{
			Properties: amqp091.Table{
				"connection_name": "gateway",
			},
		}),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
	}
	defer conn.Close()

	sc, err := microservices.NewLocalServiceConnectionContainer(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize service container: %v", err))
	}
	defer sc.Shutdown()

	// DEPENDENCIES NOW CAN BE INJECTED INTO CONCRETE API SERVICES
	fileService := file.NewLocalFileService(sc)
	authService := auth.NewLocalAuthService(sc)
	diagnoseService := diagnose.NewLocalDiagnoseService(sc)

	// Initialize Server and inject dependencies
	config := &server.FiberServerConfig{
		Host: "",
		Port: "7777",
	}

	s := server.NewFiberServer(
		config,
		fileService,
		authService,
		diagnoseService,
	)
	s.Start()
}
