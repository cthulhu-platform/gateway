package local

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/cthulhu-platform/gateway/internal/pkg"
	_ "modernc.org/sqlite"
)

type AuthRepository interface {
	GetDB() *sql.DB
	Close() error
}

type localAuthRepository struct {
	db *sql.DB
}

func NewLocalAuthRepository() (*localAuthRepository, error) {
	path := pkg.LOCAL_AUTH_REPO

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Open SQLite database connection
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	r := &localAuthRepository{
		db: db,
	}

	return r, nil
}

func (r *localAuthRepository) GetDB() *sql.DB {
	return r.db
}

func (r *localAuthRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
