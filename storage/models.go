package storage

import (
	"database/sql"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Snippet struct {
	PK            int       `json:"-"`
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	BurnAfterRead bool      `json:"burn_after_read"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type SnippetExpirationChoice struct {
	Label string
	Value string
}
type SnippetExpiration int

const (
	Never SnippetExpiration = iota
	BurnAfter
	OneHour
	OneDay
	OneWeek
	OneMonth
)

func GetSnippetExpiration(value string) SnippetExpiration {
	switch value {
	case "never":
		return Never
	case "burn_after":
		return BurnAfter
	case "one_hour":
		return OneHour
	case "one_day":
		return OneDay
	case "one_week":
		return OneWeek
	case "one_month":
		return OneMonth
	default:
		return Never
	}
}

func GetSnippetExpirationChoices() []SnippetExpirationChoice {
	return []SnippetExpirationChoice{
		{"Never", "never"},
		{"Burn After Read", "burn_after"},
		{"One Hour", "one_hour"},
		{"One Day", "one_day"},
		{"One Week", "one_week"},
		{"One Month", "one_month"},
	}
}

func (s SnippetExpiration) GetExpirationTime() *time.Time {
	switch s {
	case OneHour:
		t := time.Now().Add(time.Hour)
		return &t
	case OneDay:
		t := time.Now().Add(time.Hour * 24)
		return &t
	case OneWeek:
		t := time.Now().Add(time.Hour * 24 * 7)
		return &t
	case OneMonth:
		t := time.Now().Add(time.Hour * 24 * 30)
		return &t
	default:
		return nil
	}
}

func (s *Store) CreateSnippet(text string, expiry SnippetExpiration) (*Snippet, error) {
	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	expiresAt := sql.NullTime{}
	if expirationTime := expiry.GetExpirationTime(); expirationTime != nil {
		expiresAt.Time = *expirationTime
		expiresAt.Valid = true
	}
	burnAfterRead := false
	if expiry == BurnAfter {
		burnAfterRead = true
	}

	query := `
        INSERT INTO snippet (id, text, burn_after_read, expires_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err = s.db.client.Exec(query, id, text, burnAfterRead, expiresAt)
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
	if snippet := s.cache.client.Get(id); snippet != nil {
		return snippet, nil
	}
	query := `
		SELECT pk, id, text, burn_after_read, expires_at, created_at
		FROM snippet
		WHERE id = $1
	`
	row := s.db.client.QueryRow(query, id)
	var snippet Snippet
	var expiresAt sql.NullTime
	err := row.Scan(&snippet.PK, &snippet.ID, &snippet.Text, &snippet.BurnAfterRead, &expiresAt, &snippet.CreatedAt)
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid {
		snippet.ExpiresAt = expiresAt.Time
	}
	s.cache.client.Put(id, &snippet)
	return &snippet, nil
}
