package filemanager

import (
	"context"
	"fmt"
	"io"

	"github.com/wagslane/go-rabbitmq"
)

// FileInfo represents a stored object.
type FileInfo struct {
	FileName    string `json:"file_name"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// UploadResult is returned after an upload transaction.
type UploadResult struct {
	TransactionID string     `json:"transaction_id"`
	Success       bool       `json:"success"`
	Error         string     `json:"error,omitempty"`
	StorageID     string     `json:"storage_id,omitempty"`
	Files         []FileInfo `json:"files,omitempty"`
	TotalSize     int64      `json:"total_size,omitempty"`
}

// BucketMetadata contains objects under a storage ID.
type BucketMetadata struct {
	StorageID string     `json:"storage_id"`
	Files     []FileInfo `json:"files"`
	TotalSize int64      `json:"total_size"`
}

// DownloadResult wraps object body and metadata for streaming.
type DownloadResult struct {
	Body           io.ReadCloser
	ContentType    string
	ContentLength  int64
	DownloadedFile string
}

// UploadObject is a single file to upload.
type UploadObject struct {
	Name        string
	Size        int64
	ContentType string
	Body        io.Reader
}

type FilemanagerConnection interface {
	Upload(ctx context.Context, storageID string, objects []UploadObject) (*UploadResult, error)
	List(ctx context.Context, storageID string) (*BucketMetadata, error)
	Download(ctx context.Context, storageID, filename string) (*DownloadResult, error)
}

type rmqFilemanagerConnection struct {
	// RMQ implementation will be provided later.
}

func NewRMQFilemanagerConn(_ *rabbitmq.Conn) *rmqFilemanagerConnection {
	return &rmqFilemanagerConnection{}
}

func (c *rmqFilemanagerConnection) Upload(ctx context.Context, storageID string, objects []UploadObject) (*UploadResult, error) {
	return nil, fmt.Errorf("rmq filemanager not implemented")
}

func (c *rmqFilemanagerConnection) List(ctx context.Context, storageID string) (*BucketMetadata, error) {
	return nil, fmt.Errorf("rmq filemanager not implemented")
}

func (c *rmqFilemanagerConnection) Download(ctx context.Context, storageID, filename string) (*DownloadResult, error) {
	return nil, fmt.Errorf("rmq filemanager not implemented")
}
