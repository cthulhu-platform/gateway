package authentication

import (
	"github.com/cthulhu-platform/gateway/internal/repository/local"
)

type localAuthenticationConnection struct {
	repo local.AuthRepository
}

func NewLocalAuthConnection() (*localAuthenticationConnection, error) {
	repo, err := local.NewLocalAuthRepository()
	if err != nil {
		return nil, err
	}

	c := &localAuthenticationConnection{
		repo: repo,
	}
	return c, nil
}

func (c *localAuthenticationConnection) Close() {
	if c.repo != nil {
		c.repo.Close()
	}
}
