package implementations

import (
	"context"
	"database/sql"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository/sqlite"
)

func setupStudySessionTestDB(t *testing.T) (*sqlite.Database, func()) {
	sqlDB, err := sql.Open("sqlite3", ":memory:")

    if err != nil {
        t.Fatalf("error opening test database: %v", err)
    }

    db := &sqlite.Database{DB: sqlDB}

	// Run migrations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			words_count INTEGER DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER,
			study_activity_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS word_review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER,
			study_session_id INTEGER,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("error creating test tables: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO groups (id, name) VALUES (1, 'Test Group');
		INSERT INTO study_activities (id, name, url) VALUES (1, 'Test Activity', 'http://test.com');
	`)
	if err != nil {
		t.Fatalf("error inserting test data: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func TestStudySessionRepository_Create(t *testing.T) {
	db, cleanup := setupStudySessionTestDB(t)
	defer cleanup()

	repo := NewStudySessionRepository(db)
	ctx := context.Background()

	session := &models.StudySession{
		GroupID:         1,
		StudyActivityID: 1,
	}

	err := repo.Create(ctx, session)
	if err != nil {
		t.Errorf("error creating study session: %v", err)
	}

	if session.ID == 0 {
		t.Error("expected study session ID to be set after creation")
	}

	if session.CreatedAt.IsZero() {
		t.Error("expected created_at to be set after creation")
	}
}

func TestStudySessionRepository_GetByID(t *testing.T) {
	db, cleanup := setupStudySessionTestDB(t)
	defer cleanup()

	repo := NewStudySessionRepository(db)
	ctx := context.Background()

	// Create a test session
	session := &models.StudySession{
		GroupID:         1,
		StudyActivityID: 1,
	}
	err := repo.Create(ctx, session)
	if err != nil {
		t.Fatalf("error creating test session: %v", err)
	}

	// Test getting the session
	retrieved, err := repo.GetByID(ctx, session.ID)
	if err != nil {
		t.Errorf("error getting study session: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve study session, got nil")
	}

	if retrieved.GroupID != session.GroupID {
		t.Errorf("expected group ID %d, got %d", session.GroupID, retrieved.GroupID)
	}
}

func TestStudySessionRepository_AddReview(t *testing.T) {
	db, cleanup := setupStudySessionTestDB(t)
	defer cleanup()

	repo := NewStudySessionRepository(db)
	ctx := context.Background()

	// Create a test session
	session := &models.StudySession{
		GroupID:         1,
		StudyActivityID: 1,
	}
	err := repo.Create(ctx, session)
	if err != nil {
		t.Fatalf("error creating test session: %v", err)
	}

	// Add a review
	review := &models.WordReviewItem{
		WordID:         1,
		StudySessionID: session.ID,
		Correct:        true,
	}

	err = repo.AddReview(ctx, review)
	if err != nil {
		t.Errorf("error adding review: %v", err)
	}

	if review.ID == 0 {
		t.Error("expected review ID to be set after creation")
	}

	if review.CreatedAt.IsZero() {
		t.Error("expected created_at to be set after creation")
	}
}

// Add more tests for ListByGroup, GetSessionStats, and ListReviews methods... 