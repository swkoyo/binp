package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DBStore struct {
	client *sql.DB
}

func NewDB() (*DBStore, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))
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
            pk SERIAL PRIMARY KEY,
            id VARCHAR(50) UNIQUE NOT NULL,
            text TEXT NOT NULL,
			burn_after_read BOOLEAN NOT NULL DEFAULT FALSE,
			is_read BOOLEAN NOT NULL DEFAULT FALSE,
			expires_at TIMESTAMPTZ DEFAULT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
    `
	_, err := s.client.Exec(query)
	return err
}

func (s *DBStore) Close() error {
	return s.client.Close()
}
