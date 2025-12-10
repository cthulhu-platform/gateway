package microservices

import (
	"context"

	"github.com/cthulhu-platform/gateway/internal/microservices/authentication"
	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
	"github.com/wagslane/go-rabbitmq"
)

type ServiceConnectionContainer struct {
	Filemanager    filemanager.FilemanagerConnection
	Authentication authentication.AuthenticationConnection
}

func NewServiceConnectionContainer(ctx context.Context, conn *rabbitmq.Conn) (*ServiceConnectionContainer, error) {
	fm := filemanager.NewRMQFilemanagerConn(conn)
	auth := authentication.NewRMQAuthenticationConn(conn)

	container := &ServiceConnectionContainer{
		Filemanager:    fm,
		Authentication: auth,
	}

	return container, nil
}

func (c *ServiceConnectionContainer) Shutdown() {
}
