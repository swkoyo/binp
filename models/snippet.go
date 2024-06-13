package models

import (
	"binp/storage"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Snippet struct {
	Pk        int       `json:"-"`
	Id        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateSnippet(text string) (*Snippet, error) {
	dbStore, err := storage.GetDatabaseStore()
	if err != nil {
		return nil, err
	}

	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO snippet (id, text)
        VALUES ($1, $2)
    `
	_, err = dbStore.Client.Exec(query, id, text)
	if err != nil {
		return nil, err
	}

	snippet, err := GetSnippetByID(id)
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func GetSnippetByID(id string) (*Snippet, error) {
	dbStore, err := storage.GetDatabaseStore()
	if err != nil {
		return nil, err
	}
	query := `
		SELECT pk, id, text, created_at
		FROM snippet
		WHERE id = $1
	`
	row := dbStore.Client.QueryRow(query, id)
	var snippet Snippet
	err = row.Scan(&snippet.Pk, &snippet.Id, &snippet.Text, &snippet.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &snippet, nil
}
