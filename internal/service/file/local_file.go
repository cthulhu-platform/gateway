package file

import (
	"context"
	"errors"
	"io"
	"mime/multipart"

	"github.com/cthulhu-platform/gateway/internal/microservices"
	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
)

type localFileService struct {
	conns *microservices.ServiceConnectionContainer
}

func NewLocalFileService(conns *microservices.ServiceConnectionContainer) *localFileService {
	return &localFileService{
		conns: conns,
	}
}

func (s *localFileService) UploadFiles(ctx context.Context, files []*multipart.FileHeader) (*filemanager.UploadResult, error) {
	fm := s.filemanager()
	if fm == nil {
		return nil, errors.New("filemanager connection not configured")
}

	objects := make([]filemanager.UploadObject, 0, len(files))
	closers := make([]io.Closer, 0, len(files))
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return nil, err
		}
		closers = append(closers, file)
		objects = append(objects, filemanager.UploadObject{
			Name:        fh.Filename,
			Size:        fh.Size,
			ContentType: fh.Header.Get("Content-Type"),
			Body:        file,
		})
	}
	defer func() {
		for _, f := range closers {
			f.Close()
		}
	}()

	return fm.Upload(ctx, "", objects)
}

func (s *localFileService) DownloadFile(ctx context.Context, storageID, filename string) (*filemanager.DownloadResult, error) {
	fm := s.filemanager()
	if fm == nil {
		return nil, errors.New("filemanager connection not configured")
	}
	return fm.Download(ctx, storageID, filename)
}

func (s *localFileService) RetrieveFileBucket(ctx context.Context, storageID string) (*filemanager.BucketMetadata, error) {
	fm := s.filemanager()
	if fm == nil {
		return nil, errors.New("filemanager connection not configured")
	}
	return fm.List(ctx, storageID)
}

func (s *localFileService) filemanager() filemanager.FilemanagerConnection {
	if s == nil || s.conns == nil {
	return nil
	}
	return s.conns.Filemanager
}
