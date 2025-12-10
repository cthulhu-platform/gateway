package file

import "github.com/cthulhu-platform/gateway/internal/microservices"

type localFileService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewLocalFileService(conns *microservices.ServiceConnectionContainer) *localFileService {
	return &localFileService{
		conns: conns,
	}
}

func (s *localFileService) UploadFile() error {
	return nil
}

func (s *localFileService) DownloadFile() error {
	return nil
}

func (s *localFileService) RetrieveFileBucket() error {
	return nil
}
