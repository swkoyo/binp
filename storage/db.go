package storage

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseStore struct {
	Client *sql.DB
}

var singleton *DatabaseStore
var once sync.Once

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
		Client: db,
	}, nil
}

func GetDatabaseStore() (*DatabaseStore, error) {
	var err error
	once.Do(func() {
		singleton, err = NewDatabaseStore()
	})
	return singleton, err
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
	_, err := s.Client.Exec(query)
	return err
}
