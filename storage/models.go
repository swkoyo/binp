package storage

import (
	"binp/util"
	"database/sql"
	"fmt"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

var logger = util.GetLogger()

type Snippet struct {
	PK              int       `json:"-"`
	ID              string    `json:"id"`
	Text            string    `json:"text"`
	BurnAfterRead   bool      `json:"burn_after_read"`
	Language        string    `json:"language"`
	HighlightedCode string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type SelectOption struct {
	Label string
	Value string
}

type SnippetExpiration int

const (
	OneMinute SnippetExpiration = iota
	OneHour
	OneDay
)

var ValidLanguages = []SelectOption{
	{"Plaintext", "txt"},
	{"Bash", "bash"},
	{"CSS", "css"},
	{"Docker", "dockerfile"},
	{"Go", "go"},
	{"HTML", "html"},
	{"JavaScript", "javascript"},
	{"JSON", "json"},
	{"Lua", "lua"},
	{"Nix", "nix"},
	{"Python", "python"},
	{"Rust", "rust"},
	{"SQL", "sql"},
	{"TOML", "toml"},
	{"TypeScript", "typescript"},
	{"YAML", "yaml"},
}

var ValidExpirations = []SelectOption{
	{"One Minute", "1m"},
	{"One Hour", "1h"},
	{"One Day", "1d"},
}

func GetValidLanguages() []string {
	var langs []string
	for _, v := range ValidLanguages {
		langs = append(langs, v.Value)
	}
	return langs
}

func GetValidExpirations() []string {
	var expirations []string
	for _, v := range ValidExpirations {
		expirations = append(expirations, v.Value)
	}
	return expirations
}

func IsValidExpiration(value string) bool {
	for _, v := range ValidExpirations {
		if v.Value == value {
			return true
		}
	}
	return false
}

func IsValidLanguage(value string) bool {
	for _, v := range ValidLanguages {
		if v.Value == value {
			return true
		}
	}
	return false
}

func GetSnippetExpiration(value string) SnippetExpiration {
	switch value {
	case "1m":
		return OneMinute
	case "1h":
		return OneHour
	case "1d":
		return OneDay
	default:
		return OneMinute
	}
}

func (s SnippetExpiration) GetExpirationTime() *time.Time {
	switch s {
	case OneMinute:
		t := time.Now().UTC().Add(time.Minute)
		return &t
	case OneHour:
		t := time.Now().UTC().Add(time.Hour)
		return &t
	case OneDay:
		t := time.Now().UTC().Add(time.Hour * 24)
		return &t
	default:
		return nil
	}
}

func (s *Store) CreateSnippet(text string, burnAfterRead bool, expiry SnippetExpiration, language string) (*Snippet, error) {
	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if expirationTime := expiry.GetExpirationTime(); expirationTime != nil {
		expiresAt = expirationTime
	}

	query := `
        INSERT INTO snippet (id, text, burn_after_read, language, expires_at)
        VALUES (?, ?, ?, ?, ?)
    `
	_, err = s.db.client.Exec(query, id, text, burnAfterRead, language, expiresAt)
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
		SELECT pk, id, text, burn_after_read, language, expires_at, created_at
		FROM snippet
		WHERE id = ?
	`
	row := s.db.client.QueryRow(query, id)
	var snippet Snippet
	var expiresAt sql.NullTime
	err := row.Scan(&snippet.PK, &snippet.ID, &snippet.Text, &snippet.BurnAfterRead, &snippet.Language, &expiresAt, &snippet.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if expiresAt.Valid {
		snippet.ExpiresAt = expiresAt.Time
	}
	highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to highlight code")
		highlightedCode = snippet.Text
	}
	snippet.HighlightedCode = highlightedCode
	s.cache.client.Put(id, &snippet)
	return &snippet, nil
}

func (s *Store) UpdateSnippet(snippet *Snippet) error {
	query := `
		UPDATE snippet
		SET text = ?, burn_after_read = ?, expires_at = ?, language = ?
		WHERE id = ?
	`
	_, err := s.db.client.Exec(query, snippet.Text, snippet.BurnAfterRead, snippet.ExpiresAt, snippet.Language, snippet.ID)
	if err != nil {
		return err
	}
	highlightedCode, err := util.HighlightCode(snippet.Text, snippet.Language)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to highlight code")
		highlightedCode = snippet.Text
	}
	snippet.HighlightedCode = highlightedCode
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

func (s *Store) DeleteExpiredSnippets() (int, error) {
	ids, err := s.getExpiredSnippetIDs()
	if err != nil {
		return len(ids), err
	}
	if len(ids) == 0 {
		return len(ids), nil
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
		return len(ids), err
	}

	for _, id := range ids {
		s.cache.client.Delete(id)
	}

	return len(ids), nil
}
