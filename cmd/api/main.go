package main

import (
	"context"
	"fmt"

	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/cthulhu-platform/gateway/internal/pkg"
	"github.com/cthulhu-platform/gateway/internal/server"
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/cthulhu-platform/gateway/internal/service/file"
	"github.com/wagslane/go-rabbitmq"
)

func main() {
	// Depency Initialization (RabbitMQ conn, DB conn, SLogger)
	ctx := context.Background()

	// Create RabbitMQ connection
	connectionString := fmt.Sprintf("amqp://%s:%s@%s",
		pkg.AMPQ_HOST,
		pkg.AMPQ_PASS,
		pkg.AMPQ_HOST,
	)
	conn, err := rabbitmq.NewConn(
		connectionString,
		rabbitmq.WithConnectionOptionsLogging,
	)

	// Create DB connection and migrate
	// Create connection

	// TODO: Migration here

	sc, err := microservices.NewServiceConnectionContainer(ctx, conn)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize service container: %v", err))
	}
	defer sc.Shutdown()

	// DEPENDENCIES NOW CAN BE INJECTED INTO CONCRETE API SERVICES
	fileService := file.NewRMQFileService(sc)
	authService := auth.NewRMQAuthService(sc)

	// Initialize Server and inject dependencies
	config := &server.FiberServerConfig{
		Host: "",
		Port: "7777",
	}

	s := server.NewFiberServer(
		config,
		fileService,
		authService,
	)
	s.Start()
}
