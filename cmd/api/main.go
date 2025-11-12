package main

import (
	"context"
	"fmt"

	"github.com/cthulhu-platform/gateway/internal/server"
	"github.com/cthulhu-platform/gateway/internal/service"
)

func main() {
	// Depency Initialization (RabbitMQ conn, DB conn, SLogger)
	ctx := context.Background()

	// Initialize Server and inject dependencies
	config := &server.FiberServerConfig{
		Host: "",
		Port: "7777",
	}

	cfg := &service.ServiceContainerConfig{
		// TODO: config here
	}

	sc, err := service.NewServiceContainer(ctx, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize service container: %v", err))
	}
	defer sc.Shutdown()

	s := server.NewFiberServer(config, sc)
	s.Start()
}
