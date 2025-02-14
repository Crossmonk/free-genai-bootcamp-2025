package test

import (
	"context"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/service"
)


func TestStudyActivityRepository_Create(t *testing.T) {
	repo := service.NewMockStudyActivityRepository()
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
	repo := service.NewMockStudyActivityRepository()
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
}

func TestStudyActivityRepository_List(t *testing.T) {
	repo := service.NewMockStudyActivityRepository()
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
	retrievedActivities, err := repo.List(ctx)
	if err != nil {
		t.Errorf("error listing study activities: %v", err)
	}

	if len(retrievedActivities) != len(activities) {
		t.Errorf("expected %d activities, got %d", len(activities), len(retrievedActivities))
	}
}

func TestStudyActivityRepository_Update(t *testing.T) {
	repo := service.NewMockStudyActivityRepository()
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

	// Update the activity
	activity.Name = "Updated Flashcards"
	err = repo.Update(ctx, activity)
	if err != nil {
		t.Errorf("error updating study activity: %v", err)
	}

	// Verify the update
	updatedActivity, err := repo.GetByID(ctx, activity.ID)
	if err != nil {
		t.Errorf("error getting updated activity: %v", err)
	}

	if updatedActivity.Name != "Updated Flashcards" {
		t.Errorf("expected updated name to be 'Updated Flashcards', got %s", updatedActivity.Name)
	}
}

func TestStudyActivityRepository_Delete(t *testing.T) {
	repo := service.NewMockStudyActivityRepository()
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

	// Delete the activity
	err = repo.Delete(ctx, activity.ID)
	if err != nil {
		t.Errorf("error deleting study activity: %v", err)
	}

	// Verify deletion
	deletedActivity, err := repo.GetByID(ctx, activity.ID)
	if err != nil {
		t.Errorf("error getting deleted activity: %v", err)
	}

	if deletedActivity != nil {
		t.Fatal("expected study activity to be deleted, got non-nil")
	}
} 