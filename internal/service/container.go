package service

import "context"

type ServiceContainerConfig struct{}

type ServiceContainer struct {
	FileManager FileManagerService
}

// Creates a new service container and starts it up
// Returns error if startup failed
func NewServiceContainer(ctx context.Context, cfg *ServiceContainerConfig) (*ServiceContainer, error) {
	// Context is going to have a timeout
	// Establish rabbitmq connection

	fm := NewRMQFileManagerService()

	sc := &ServiceContainer{
		FileManager: fm,
	}

	return sc, nil
}

// Called in main after fiber app shuts down
func (sc *ServiceContainer) Shutdown() error {
	// TODO: Shutdown connections to repo or rabbitmq here
	return nil
}
