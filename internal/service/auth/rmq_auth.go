package auth

import "github.com/cthulhu-platform/gateway/internal/microservices"

type rmqAuthService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewRMQAuthService(conns *microservices.ServiceConnectionContainer) *rmqAuthService {
	return &rmqAuthService{
		conns: conns,
	}
}

func (s *rmqAuthService) Validate() {}
