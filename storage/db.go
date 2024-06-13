package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseStore struct {
	client *sql.DB
}

func NewDatabaseStore() (*DatabaseStore, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

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

	return &DatabaseStore{
		client: db,
	}, nil
}

func (s *DatabaseStore) Init() error {
	return s.createSnippetTable()
}

func (s *DatabaseStore) createSnippetTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS snippet (
            pk SERIAL PRIMARY KEY,
            id VARCHAR(50) UNIQUE NOT NULL,
            text TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        )
    `
	_, err := s.client.Exec(query)
	return err
}
