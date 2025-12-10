package file

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
)

type rmqFileService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewRMQFileService(conns *microservices.ServiceConnectionContainer) *rmqFileService {
	return &rmqFileService{
		conns: conns,
	}
}

func (s *rmqFileService) UploadFiles(ctx context.Context, files []*multipart.FileHeader) (*filemanager.UploadResult, error) {
	// RMQ implementation will be added later.
	return nil, errors.New("rmq file service not implemented")
}

func (s *rmqFileService) DownloadFile(ctx context.Context, storageID, filename string) (*filemanager.DownloadResult, error) {
	return nil, errors.New("rmq file service not implemented")
}

func (s *rmqFileService) RetrieveFileBucket(ctx context.Context, storageID string) (*filemanager.BucketMetadata, error) {
	return nil, errors.New("rmq file service not implemented")
}

func (s *rmqFileService) filemanager() filemanager.FilemanagerConnection {
	if s == nil || s.conns == nil {
	return nil
	}
	return s.conns.Filemanager
}
