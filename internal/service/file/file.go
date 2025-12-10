package file

import (
	"context"

	"mime/multipart"

	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
)

type FileService interface {
	UploadFiles(ctx context.Context, files []*multipart.FileHeader, userID *string) (*filemanager.UploadResult, error)
	DownloadFile(ctx context.Context, storageID, stringID string) (*filemanager.DownloadResult, error)
	RetrieveFileBucket(ctx context.Context, storageID string) (*filemanager.BucketMetadata, error)
}
