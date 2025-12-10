package filemanager

import (
	"github.com/cthulhu-platform/gateway/internal/repository/local"
)

type localFilemanagerConnection struct {
	repo local.FileRepository
}

func NewLocalFilemanagerConnection() (*localFilemanagerConnection, error) {
	repo, err := local.NewLocalFileRepository()
	if err != nil {
		return nil, err
	}

	c := &localFilemanagerConnection{
		repo: repo,
	}
	return c, nil
}

func (c *localFilemanagerConnection) Close() {
	if c.repo != nil {
		c.repo.Close()
	}
}
