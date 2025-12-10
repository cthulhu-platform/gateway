package file

import "github.com/cthulhu-platform/gateway/internal/microservices"

type rmqFileService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewRMQFileService(conns *microservices.ServiceConnectionContainer) *rmqFileService {
	return &rmqFileService{
		conns: conns,
	}
}

func (s *rmqFileService) UploadFile() error {
	return nil
}

func (s *rmqFileService) DownloadFile() error {
	return nil
}

func (s *rmqFileService) RetrieveFileBucket() error {
	return nil
}
