package implementations

import (
	"context"
	"database/sql"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository/sqlite"
)

func setupStudyActivityTestDB(t *testing.T) (*sqlite.Database, func()) {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("error opening test database: %v", err)
    }

    db := &sqlite.Database{DB: sqlDB} 

	// Run migrations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			study_activity_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("error creating test tables: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func TestStudyActivityRepository_Create(t *testing.T) {
	db, cleanup := setupStudyActivityTestDB(t)
	defer cleanup()

	repo := NewStudyActivityRepository(db)
	ctx := context.Background()

	activity := &models.StudyActivity{
		Name: "Flashcards",
		URL:  "https://example.com/flashcards",
	}

	err := repo.Create(ctx, activity)
	if err != nil {
		t.Errorf("error creating study activity: %v", err)
	}

	if activity.ID == 0 {
		t.Error("expected study activity ID to be set after creation")
	}
}

func TestStudyActivityRepository_GetByID(t *testing.T) {
	db, cleanup := setupStudyActivityTestDB(t)
	defer cleanup()

	repo := NewStudyActivityRepository(db)
	ctx := context.Background()

	// Create a test activity
	activity := &models.StudyActivity{
		Name: "Flashcards",
		URL:  "https://example.com/flashcards",
	}
	err := repo.Create(ctx, activity)
	if err != nil {
		t.Fatalf("error creating test activity: %v", err)
	}

	// Test getting the activity
	retrieved, err := repo.GetByID(ctx, activity.ID)
	if err != nil {
		t.Errorf("error getting study activity: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve study activity, got nil")
	}

	if retrieved.Name != activity.Name {
		t.Errorf("expected name %s, got %s", activity.Name, retrieved.Name)
	}

	if retrieved.URL != activity.URL {
		t.Errorf("expected URL %s, got %s", activity.URL, retrieved.URL)
	}
}

func TestStudyActivityRepository_List(t *testing.T) {
	db, cleanup := setupStudyActivityTestDB(t)
	defer cleanup()

	repo := NewStudyActivityRepository(db)
	ctx := context.Background()

	// Create test activities
	activities := []*models.StudyActivity{
		{Name: "Flashcards", URL: "https://example.com/flashcards"},
		{Name: "Quiz", URL: "https://example.com/quiz"},
	}

	for _, activity := range activities {
		err := repo.Create(ctx, activity)
		if err != nil {
			t.Fatalf("error creating test activity: %v", err)
		}
	}

	// Test listing activities
	retrieved, err := repo.List(ctx)
	if err != nil {
		t.Errorf("error listing study activities: %v", err)
	}

	if len(retrieved) != len(activities) {
		t.Errorf("expected %d activities, got %d", len(activities), len(retrieved))
	}
}

// Add more tests for Update and Delete methods... 