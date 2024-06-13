package storage

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Snippet struct {
	PK        int       `json:"-"`
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) createSnippetTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS snippet (
            pk SERIAL PRIMARY KEY,
            id VARCHAR(50) UNIQUE NOT NULL,
            text TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        )
    `
	_, err := s.db.client.Exec(query)
	return err
}

func (s *Store) CreateSnippet(text string) (*Snippet, error) {
	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO snippet (id, text)
        VALUES ($1, $2)
    `
	_, err = s.db.client.Exec(query, id, text)
	if err != nil {
		return nil, err
	}

	snippet, err := s.GetSnippetByID(id)
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *Store) GetSnippetByID(id string) (*Snippet, error) {
	query := `
		SELECT pk, id, text, created_at
		FROM snippet
		WHERE id = $1
	`
	row := s.db.client.QueryRow(query, id)
	var snippet Snippet
	err := row.Scan(&snippet.PK, &snippet.ID, &snippet.Text, &snippet.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &snippet, nil
}
