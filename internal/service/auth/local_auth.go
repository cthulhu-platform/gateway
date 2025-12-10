package auth

import "github.com/cthulhu-platform/gateway/internal/microservices"

type localAuthService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewLocalAuthService(conns *microservices.ServiceConnectionContainer) *localAuthService {
	return &localAuthService{
		conns: conns,
	}
}

func (s *localAuthService) Validate() {
}
