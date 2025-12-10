package file

import (
	"context"

	"mime/multipart"

	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
)

type FileService interface {
	UploadFiles(ctx context.Context, files []*multipart.FileHeader) (*filemanager.UploadResult, error)
	DownloadFile(ctx context.Context, storageID, filename string) (*filemanager.DownloadResult, error)
	RetrieveFileBucket(ctx context.Context, storageID string) (*filemanager.BucketMetadata, error)
}
