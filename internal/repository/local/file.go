package local

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"github.com/cthulhu-platform/gateway/internal/pkg"
)

type FileRepository interface {
	GetDB() *sql.DB
	Close() error
}

type localFileRepository struct {
	db *sql.DB
}

func NewLocalFileRepository() (*localFileRepository, error) {
	path := pkg.LOCAL_FILE_REPO

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

	r := &localFileRepository{
		db: db,
	}

	return r, nil
}

func (r *localFileRepository) GetDB() *sql.DB {
	return r.db
}

func (r *localFileRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
