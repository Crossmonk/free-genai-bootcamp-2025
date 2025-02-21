package test

import (
	"context"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/service"
)

func TestGroupRepository_Create(t *testing.T) {
	repo := service.NewMockGroupRepository()
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

// Add more tests for GetByID, List, Update, and Delete methods... 

func TestGroupRepository_GetByID(t *testing.T) {
	repo := service.NewMockGroupRepository()
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

func TestGroupRepository_List(t *testing.T) {
	repo := service.NewMockGroupRepository()
	ctx := context.Background()

	// Create test groups
	groups := []*models.Group{
		{Name: "Basic Verbs"},
		{Name: "Advanced Verbs"},
	}

	for _, group := range groups {
		err := repo.Create(ctx, group)
		if err != nil {
			t.Fatalf("error creating test group: %v", err)
		}
	}

	// Test listing groups
	retrievedGroups, total, err := repo.List(ctx, 1, 10, "name", "asc")
	if err != nil {
		t.Errorf("error listing groups: %v", err)
	}

	if len(retrievedGroups) != len(groups) {
		t.Errorf("expected %d groups, got %d", len(groups), len(retrievedGroups))
	}

	if total != len(groups) {
		t.Errorf("expected total %d, got %d", len(groups), total)
	}
}

func TestGroupRepository_Update(t *testing.T) {
	repo := service.NewMockGroupRepository()
	ctx := context.Background()

	// Create a test group
	group := &models.Group{
		Name: "Basic Verbs",
	}
	err := repo.Create(ctx, group)
	if err != nil {
		t.Fatalf("error creating test group: %v", err)
	}

	// Update the group
	group.Name = "Updated Basic Verbs"
	err = repo.Update(ctx, group)
	if err != nil {
		t.Errorf("error updating group: %v", err)
	}

	// Verify the update
	updatedGroup, err := repo.GetByID(ctx, group.ID)
	if err != nil {
		t.Errorf("error getting updated group: %v", err)
	}

	if updatedGroup.Name != "Updated Basic Verbs" {
		t.Errorf("expected updated name to be 'Updated Basic Verbs', got %s", updatedGroup.Name)
	}
}

func TestGroupRepository_Delete(t *testing.T) {
	// db, cleanup := setupGroupTestDB(t)
	// defer cleanup()

	repo := service.NewMockGroupRepository()
	ctx := context.Background()

	// Create a test group
	group := &models.Group{
		Name: "Basic Verbs",
	}
	err := repo.Create(ctx, group)
	if err != nil {
		t.Fatalf("error creating test group: %v", err)
	}

	// Delete the group
	err = repo.Delete(ctx, group.ID)
	if err != nil {
		t.Errorf("error deleting group: %v", err)
	}

	// Verify deletion
	deletedGroup, err := repo.GetByID(ctx, group.ID)
	if err != nil {
		t.Errorf("error getting deleted group: %v", err)
	}

	if deletedGroup != nil {
		t.Fatal("expected group to be deleted, got non-nil")
	}
} 