package implementations

import (
	"context"
	"database/sql"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository/sqlite"
)

func setupGroupTestDB(t *testing.T) (*sqlite.Database, func()) {
	sqlDB, err := sql.Open("sqlite3", ":memory:")

    if err != nil {
        t.Fatalf("error opening test database: %v", err)
    }

    db := &sqlite.Database{DB: sqlDB}

	if err != nil {
		t.Fatalf("error opening test database: %v", err)
	}

	// Run migrations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			words_count INTEGER DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kanji TEXT NOT NULL,
			romaji TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT NOT NULL CHECK (json_valid(parts))
		);

		CREATE TABLE IF NOT EXISTS word_groups (
			word_id INTEGER,
			group_id INTEGER,
			PRIMARY KEY (word_id, group_id),
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER,
			study_activity_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS word_review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER,
			study_session_id INTEGER,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("error creating test tables: %v", err)
	}

	return db, func() {
		sqlDB.Close()
	}
}

func TestGroupRepository_Create(t *testing.T) {
	db, cleanup := setupGroupTestDB(t)
	defer cleanup()

	repo := NewGroupRepository(db)
	ctx := context.Background()

	group := &models.Group{
		Name: "Basic Verbs",
	}

	err := repo.Create(ctx, group)
	if err != nil {
		t.Errorf("error creating group: %v", err)
	}

	if group.ID == 0 {
		t.Error("expected group ID to be set after creation")
	}
}

func TestGroupRepository_GetByID(t *testing.T) {
	db, cleanup := setupGroupTestDB(t)
	defer cleanup()

	repo := NewGroupRepository(db)
	ctx := context.Background()

	// Create a test group
	group := &models.Group{
		Name: "Basic Verbs",
	}
	err := repo.Create(ctx, group)
	if err != nil {
		t.Fatalf("error creating test group: %v", err)
	}

	// Test getting the group
	retrieved, err := repo.GetByID(ctx, group.ID)
	if err != nil {
		t.Errorf("error getting group: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve group, got nil")
	}

	if retrieved.Name != group.Name {
		t.Errorf("expected name %s, got %s", group.Name, retrieved.Name)
	}
}

// Add more tests for other methods... 