package storage

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBStore struct {
	client *sql.DB
}

func NewDB() (*DBStore, error) {
	dbPath := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DBStore{
		client: db,
	}, nil
}

func (s *DBStore) Init() error {
	query := `
        CREATE TABLE IF NOT EXISTS snippet (
            pk INTEGER PRIMARY KEY AUTOINCREMENT,
            id TEXT UNIQUE NOT NULL,
            text TEXT NOT NULL,
			burn_after_read INTEGER NOT NULL DEFAULT 0,
			is_read INTEGER NOT NULL DEFAULT 0,
			expires_at DATETIME DEFAULT NULL,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `
	_, err := s.client.Exec(query)
	return err
}

func (s *DBStore) Close() error {
	return s.client.Close()
}
