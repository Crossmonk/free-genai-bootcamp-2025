package implementations

import (
	"context"
	"database/sql"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sqlite.Database, func()) {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("error opening test database: %v", err)
    }

    db := &sqlite.Database{DB: sqlDB}  // Wrap the sql.DB

	// Run migrations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kanji TEXT NOT NULL,
			romaji TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT NOT NULL CHECK (json_valid(parts))
		);

		CREATE TABLE IF NOT EXISTS word_review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER,
			study_session_id INTEGER,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("error creating test tables: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func TestWordRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewWordRepository(db)
	ctx := context.Background()

	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts: map[string]any{
			"verb_type": "ru-verb",
			"topic":     "food",
		},
	}

	err := repo.Create(ctx, word)
	if err != nil {
		t.Errorf("error creating word: %v", err)
	}

	if word.ID == 0 {
		t.Error("expected word ID to be set after creation")
	}
}

func TestWordRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewWordRepository(db)
	ctx := context.Background()

	// Create a test word
	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts: map[string]any{
			"verb_type": "ru-verb",
			"topic":     "food",
		},
	}
	err := repo.Create(ctx, word)
	if err != nil {
		t.Fatalf("error creating test word: %v", err)
	}

	// Test getting the word
	retrieved, err := repo.GetByID(ctx, word.ID)
	if err != nil {
		t.Errorf("error getting word: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve word, got nil")
	}

	if retrieved.Kanji != word.Kanji {
		t.Errorf("expected kanji %s, got %s", word.Kanji, retrieved.Kanji)
	}
}

// Add more tests for List, Update, Delete, and GetStats methods... 