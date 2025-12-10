package repository

type libsqlSessionRepository struct{}

func NewLibsqlSessionRepository() (*libsqlSessionRepository, error) {
	rep := &libsqlSessionRepository{}
	return rep, nil
}
