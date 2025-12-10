package file

import (
	"context"

	"mime/multipart"

	"github.com/cthulhu-platform/gateway/internal/microservices/filemanager"
)

type AdminInfo struct {
	UserID    string  `json:"user_id"`
	Email     string  `json:"email"`
	Username  *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsOwner   bool    `json:"is_owner"`
	CreatedAt int64   `json:"created_at"`
}

type BucketAdminsResponse struct {
	BucketID string      `json:"bucket_id"`
	Owner    *AdminInfo  `json:"owner"`
	Admins   []AdminInfo `json:"admins"`
}

type FileService interface {
	UploadFiles(ctx context.Context, files []*multipart.FileHeader, userID *string) (*filemanager.UploadResult, error)
	DownloadFile(ctx context.Context, storageID, stringID string) (*filemanager.DownloadResult, error)
	RetrieveFileBucket(ctx context.Context, storageID string) (*filemanager.BucketMetadata, error)
	GetBucketAdmins(ctx context.Context, bucketID string) (*BucketAdminsResponse, error)
}
