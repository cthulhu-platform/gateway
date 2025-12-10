package diagnose

import "github.com/cthulhu-platform/gateway/internal/microservices"

type localDiagnoseService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewLocalDiagnoseService(conns *microservices.ServiceConnectionContainer) *localDiagnoseService {
	return &localDiagnoseService{
		conns: conns,
	}
}

func (s *localDiagnoseService) ServiceFanoutTest() (string, error) {
	return "", nil
}
