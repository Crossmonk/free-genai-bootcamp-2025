package service

import (
	"context"
	"testing"

	"backend-go/internal/domain/models"
)

func newMockStudyActivityRepository() *mockStudyActivityRepository {
	return &mockStudyActivityRepository{
		activities: make(map[int64]*models.StudyActivity),
	}
}

func newMockStudySessionRepository() *mockStudySessionRepository {
	return &mockStudySessionRepository{
		// Initialize as needed
	}
}

func TestStudyActivityService_CreateActivity(t *testing.T) {
	repo := newMockStudyActivityRepository()
	sessionRepo := newMockStudySessionRepository()
	service := NewStudyActivityService(repo, sessionRepo)
	ctx := context.Background()

	tests := []struct {
		name     string
		activity *models.StudyActivity
		wantErr  bool
	}{
		{
			name: "valid activity",
			activity: &models.StudyActivity{
				Name: "Flashcards",
				URL:  "https://example.com/flashcards",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			activity: &models.StudyActivity{
				URL: "https://example.com/flashcards",
			},
			wantErr: true,
		},
		{
			name: "missing URL",
			activity: &models.StudyActivity{
				Name: "Flashcards",
			},
			wantErr: true,
		},
		{
			name: "invalid URL",
			activity: &models.StudyActivity{
				Name: "Flashcards",
				URL:  "not-a-url",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateActivity(ctx, tt.activity)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.activity.ID == 0 {
				t.Error("CreateActivity() did not set ID for valid activity")
			}
		})
	}
}

func TestStudyActivityService_GetActivity(t *testing.T) {
	repo := newMockStudyActivityRepository()
	sessionRepo := newMockStudySessionRepository()
	service := NewStudyActivityService(repo, sessionRepo)
	ctx := context.Background()

	// Create a test activity
	activity := &models.StudyActivity{
		Name: "Flashcards",
		URL:  "https://example.com/flashcards",
	}
	err := service.CreateActivity(ctx, activity)
	if err != nil {
		t.Fatalf("Failed to create test activity: %v", err)
	}

	// Test getting the activity
	retrieved, err := service.GetActivity(ctx, activity.ID)
	if err != nil {
		t.Errorf("GetActivity() error = %v", err)
		return
	}

	if retrieved.ID != activity.ID {
		t.Errorf("GetActivity() got ID = %v, want %v", retrieved.ID, activity.ID)
	}

	if retrieved.Name != activity.Name {
		t.Errorf("GetActivity() got Name = %v, want %v", retrieved.Name, activity.Name)
	}

	if retrieved.URL != activity.URL {
		t.Errorf("GetActivity() got URL = %v, want %v", retrieved.URL, activity.URL)
	}
}

func TestStudyActivityService_ListActivities(t *testing.T) {
	repo := newMockStudyActivityRepository()
	sessionRepo := newMockStudySessionRepository()
	service := NewStudyActivityService(repo, sessionRepo)
	ctx := context.Background()

	// Create test activities
	activities := []*models.StudyActivity{
		{Name: "Flashcards", URL: "https://example.com/flashcards"},
		{Name: "Quiz", URL: "https://example.com/quiz"},
	}

	for _, activity := range activities {
		err := service.CreateActivity(ctx, activity)
		if err != nil {
			t.Fatalf("Failed to create test activity: %v", err)
		}
	}

	// Test listing activities
	retrieved, err := service.ListActivities(ctx)
	if err != nil {
		t.Errorf("ListActivities() error = %v", err)
		return
	}

	if len(retrieved) != len(activities) {
		t.Errorf("ListActivities() got %v activities, want %v", len(retrieved), len(activities))
	}
} 