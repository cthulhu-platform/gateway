package microservices

import (
	"context"

	"github.com/cthulhu-platform/gateway/internal/microservices/authentication"
	"github.com/cthulhu-platform/gateway/internal/microservices/diagnose"
	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
	"github.com/wagslane/go-rabbitmq"
)

type ServiceConnectionContainer struct {
	Filemanager    filemanager.FilemanagerConnection
	Authentication authentication.AuthenticationConnection
	Diagnose       diagnose.DiagnoseConnection
}

func NewLocalServiceConnectionContainer(ctx context.Context) (*ServiceConnectionContainer, error) {
	fm, err := filemanager.NewLocalFilemanagerConnection()
	if err != nil {
		return nil, err
	}
	auth, err := authentication.NewLocalAuthConnection()
	if err != nil {
		return nil, err
	}
	diag, err := diagnose.NewLocalDiagnoseConnection()
	if err != nil {
		return nil, err
	}

	container := &ServiceConnectionContainer{
		Filemanager:    fm,
		Authentication: auth,
		Diagnose:       diag,
	}

	return container, nil
}

func NewServiceConnectionContainer(ctx context.Context, conn *rabbitmq.Conn) (*ServiceConnectionContainer, error) {
	fm := filemanager.NewRMQFilemanagerConn(conn)
	auth, err := authentication.NewRMQAuthenticationConn(conn)
	if err != nil {
		return nil, err
	}
	diag, err := diagnose.NewRMQDiagnoseConn(conn)
	if err != nil {
		return nil, err
	}

	container := &ServiceConnectionContainer{
		Filemanager:    fm,
		Authentication: auth,
		Diagnose:       diag,
	}

	return container, nil
}

func (c *ServiceConnectionContainer) Shutdown() {
	if authConn, ok := c.Authentication.(interface{ Close() }); ok {
		authConn.Close()
	}
	if diagConn, ok := c.Diagnose.(interface{ Close() }); ok {
		diagConn.Close()
	}
}
