package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

type Snippet struct {
	PK            int       `json:"-"`
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	BurnAfterRead bool      `json:"burn_after_read"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type SnippetExpirationChoice struct {
	Label string
	Value string
}
type SnippetExpiration int

const (
	OneHour SnippetExpiration = iota
	OneDay
	OneWeek
	OneMonth
)

func GetSnippetExpiration(value string) SnippetExpiration {
	switch value {
	case "one_hour":
		return OneHour
	case "one_day":
		return OneDay
	case "one_week":
		return OneWeek
	case "one_month":
		return OneMonth
	default:
		return OneHour
	}
}

func GetSnippetExpirationChoices() []SnippetExpirationChoice {
	return []SnippetExpirationChoice{
		{"One Hour", "one_hour"},
		{"One Day", "one_day"},
		{"One Week", "one_week"},
		{"One Month", "one_month"},
	}
}

func (s SnippetExpiration) GetExpirationTime() *time.Time {
	switch s {
	case OneHour:
		t := time.Now().UTC().Add(time.Hour)
		return &t
	case OneDay:
		t := time.Now().UTC().Add(time.Hour * 24)
		return &t
	case OneWeek:
		t := time.Now().UTC().Add(time.Hour * 24 * 7)
		return &t
	case OneMonth:
		t := time.Now().UTC().Add(time.Hour * 24 * 30)
		return &t
	default:
		return nil
	}
}

func (s *Store) CreateSnippet(text string, burnAfterRead bool, expiry SnippetExpiration) (*Snippet, error) {
	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if expirationTime := expiry.GetExpirationTime(); expirationTime != nil {
		expiresAt = expirationTime
	}

	query := `
        INSERT INTO snippet (id, text, burn_after_read, expires_at)
        VALUES (?, ?, ?, ?)
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
		SELECT pk, id, text, burn_after_read, is_read, expires_at, created_at
		FROM snippet
		WHERE id = ?
	`
	row := s.db.client.QueryRow(query, id)
	var snippet Snippet
	var expiresAt sql.NullTime
	err := row.Scan(&snippet.PK, &snippet.ID, &snippet.Text, &snippet.BurnAfterRead, &snippet.IsRead, &expiresAt, &snippet.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if expiresAt.Valid {
		snippet.ExpiresAt = expiresAt.Time
		if snippet.ExpiresAt.Before(time.Now().UTC()) {
			s.cache.client.Delete(id)
			return nil, nil
		}
	}
	s.cache.client.Put(id, &snippet)
	return &snippet, nil
}

func (s *Store) UpdateSnippet(snippet *Snippet) error {
	query := `
		UPDATE snippet
		SET text = ?, burn_after_read = ?, expires_at = ?, is_read = ?
		WHERE id = ?
	`
	_, err := s.db.client.Exec(query, snippet.Text, snippet.BurnAfterRead, snippet.ExpiresAt, snippet.IsRead, snippet.ID)
	if err != nil {
		return err
	}
	s.cache.client.Put(snippet.ID, snippet)
	return nil
}

func (s *Store) DeleteSnippet(id string) error {
	query := `
		DELETE FROM snippet
		WHERE id = ?
	`
	_, err := s.db.client.Exec(query, id)
	if err != nil {
		return err
	}
	s.cache.client.Delete(id)
	return nil
}

func (s *Store) getExpiredSnippetIDs() ([]string, error) {
	query := `
		SELECT id
		FROM snippet
		WHERE expires_at <= datetime('now')
	`
	rows, err := s.db.client.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *Store) DeleteExpiredSnippets() error {
	ids, err := s.getExpiredSnippetIDs()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
		DELETE FROM snippet
		WHERE id IN (%s)
	`, strings.Join(placeholders, ", "))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = s.db.client.Exec(query, args...)
	if err != nil {
		return err
	}

	for _, id := range ids {
		s.cache.client.Delete(id)
	}

	return nil
}
