package storage

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestStore(t *testing.T) *Store {
	os.Setenv("DB_PATH", ":memory:")
	dbStore, err := NewDB()
	assert.NoError(t, err)

	err = dbStore.Init()
	assert.NoError(t, err)

	cacheStore := NewCache()

	return &Store{
		db:    dbStore,
		cache: cacheStore,
	}
}

func TestCreateSnippet(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()

	testCases := []struct {
		name             string
		expiry           SnippetExpiration
		burnAfterRead    bool
		expectedDuration time.Duration
	}{
		{"BurnAfter", OneHour, true, time.Hour},
		{"OneHour", OneHour, false, time.Hour},
		{"OneDay", OneDay, false, (24 * time.Hour)},
		{"OneWeek", OneWeek, false, (7 * 24 * time.Hour)},
		{"OneMonth", OneMonth, false, (30 * 24 * time.Hour)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			text := "Test snippet for " + tc.name
			snippet, err := store.CreateSnippet(text, tc.burnAfterRead, tc.expiry)
			if err != nil {
				t.Fatalf("Failed to create snippet: %v", err)
			}

			if snippet.ID == "" {
				t.Errorf("Expected snippet ID to be non-empty, got %q", snippet.ID)
			}

			if snippet.Text != text {
				t.Errorf("Expected snippet text to be %q, got %q", text, snippet.Text)
			}

			if snippet.BurnAfterRead != tc.burnAfterRead {
				t.Errorf("Expected snippet burnAfterRead to be %v, got %v", tc.burnAfterRead, snippet.BurnAfterRead)
			}

			if snippet.ExpiresAt.Location() != time.UTC {
				t.Errorf("Expected snippet expiration time to be in UTC, got %v", snippet.ExpiresAt.Location())
			}

			expectedTime := time.Now().UTC().Add(tc.expectedDuration)
			if snippet.ExpiresAt.Sub(expectedTime) > time.Second {
				t.Errorf("Expected snippet to expire at %v, got %v", expectedTime, snippet.ExpiresAt)
			}
		})
	}
}

func TestGetSnippetByID(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()

	createdSnippet, err := store.CreateSnippet("Test snippet", false, OneHour)
	assert.NoError(t, err)

	retrievedSnippet, err := store.GetSnippetByID(createdSnippet.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedSnippet)
	assert.Equal(t, createdSnippet.ID, retrievedSnippet.ID)
	assert.Equal(t, createdSnippet.Text, retrievedSnippet.Text)
	assert.Equal(t, createdSnippet.BurnAfterRead, retrievedSnippet.BurnAfterRead)
	assert.Equal(t, createdSnippet.ExpiresAt.Location(), time.UTC)
	assert.Equal(t, createdSnippet.ExpiresAt, retrievedSnippet.ExpiresAt)

	cachedSnippet := store.cache.client.Get(createdSnippet.ID)
	assert.NotNil(t, cachedSnippet)
	assert.Equal(t, createdSnippet.ID, cachedSnippet.ID)
	assert.Equal(t, createdSnippet.Text, cachedSnippet.Text)
	assert.Equal(t, createdSnippet.BurnAfterRead, cachedSnippet.BurnAfterRead)
	assert.Equal(t, createdSnippet.ExpiresAt.Location(), time.UTC)
	assert.Equal(t, createdSnippet.ExpiresAt, cachedSnippet.ExpiresAt)

	nonExistentSnippet, err := store.GetSnippetByID("non-existent-id")
	assert.NoError(t, err)
	assert.Nil(t, nonExistentSnippet)
}

func TestDeleteSnippet(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()

	snippet, err := store.CreateSnippet("Test snippet", false, OneHour)
	assert.NoError(t, err)

	store.cache.client.Put(snippet.ID, snippet)

	err = store.DeleteSnippet(snippet.ID)
	assert.NoError(t, err)

	deletedSnippet, err := store.GetSnippetByID(snippet.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedSnippet)

	cachedSnippet := store.cache.client.Get(snippet.ID)
	assert.Nil(t, cachedSnippet)
}

func TestUpdateSnippet(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()

	snippet, err := store.CreateSnippet("Test snippet", false, OneHour)
	assert.NoError(t, err)
	assert.False(t, snippet.IsRead)

	store.cache.client.Put(snippet.ID, snippet)

	snippet.IsRead = true
	err = store.UpdateSnippet(snippet)
	assert.NoError(t, err)

	cachedSnippet := store.cache.client.Get(snippet.ID)
	assert.True(t, cachedSnippet.IsRead)

	updatedSnippet, err := store.GetSnippetByID(snippet.ID)
	assert.NoError(t, err)
	assert.True(t, updatedSnippet.IsRead)
}

func TestGetExpiredSnippetIDs(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()

	expiredSnippet, err := store.CreateSnippet("Expired snippet", false, OneHour)
	assert.NoError(t, err)

	_, err = store.db.client.Exec(
		"UPDATE snippet SET expires_at = ? WHERE id = ?",
		time.Now().UTC().Add(-time.Hour),
		expiredSnippet.ID,
	)
	assert.NoError(t, err)

	_, err = store.CreateSnippet("Valid snippet", false, OneDay)
	assert.NoError(t, err)

	ids, err := store.getExpiredSnippetIDs()
	assert.NoError(t, err)

	assert.Len(t, ids, 1)
	assert.Equal(t, expiredSnippet.ID, ids[0])
}

func TestDeleteExpiredSnippets(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()

	expiredSnippet, err := store.CreateSnippet("Expired snippet", false, OneHour)
	assert.NoError(t, err)

	_, err = store.db.client.Exec(
		"UPDATE snippet SET expires_at = ? WHERE id = ?",
		time.Now().UTC().Add(-time.Hour),
		expiredSnippet.ID,
	)
	assert.NoError(t, err)

	store.cache.client.Put(expiredSnippet.ID, expiredSnippet)

	validSnippet, err := store.CreateSnippet("Valid snippet", false, OneDay)
	assert.NoError(t, err)

	err = store.DeleteExpiredSnippets()
	assert.NoError(t, err)

	deletedSnippet, err := store.GetSnippetByID(expiredSnippet.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedSnippet)

	cachedDeletedSnippet := store.cache.client.Get(expiredSnippet.ID)
	assert.Nil(t, cachedDeletedSnippet)

	existingSnippet, err := store.GetSnippetByID(validSnippet.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existingSnippet)
}
