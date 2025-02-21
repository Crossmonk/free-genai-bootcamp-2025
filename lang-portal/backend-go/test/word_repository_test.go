package test

import (
	"context"
	"testing"

	"backend-go/internal/domain/models"
	"backend-go/internal/service"
)

func TestWordRepository_Create(t *testing.T) {
	repo := service.NewMockWordRepository()
	ctx := context.Background()

	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts:   map[string]any{"verb_type": "ru-verb"},
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
	repo := service.NewMockWordRepository()
	ctx := context.Background()

	// Create a test word
	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts:   map[string]any{"verb_type": "ru-verb"},
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

func TestWordRepository_List(t *testing.T) {
	repo := service.NewMockWordRepository()
	ctx := context.Background()

	// Create test words
	words := []*models.Word{
		{Kanji: "食べる", Romaji: "taberu", English: "to eat", Parts: map[string]any{"verb_type": "ru-verb"}},
		{Kanji: "飲む", Romaji: "nomu", English: "to drink", Parts: map[string]any{"verb_type": "ru-verb"}},
	}

	for _, word := range words {
		err := repo.Create(ctx, word)
		if err != nil {
			t.Fatalf("error creating test word: %v", err)
		}
	}

	// Test listing words
	retrievedWords, total, err := repo.List(ctx, 1, 10, "kanji", "asc")
	if err != nil {
		t.Errorf("error listing words: %v", err)
	}

	if len(retrievedWords) != len(words) {
		t.Errorf("expected %d words, got %d", len(words), len(retrievedWords))
	}

	if total != len(words) {
		t.Errorf("expected total %d, got %d", len(words), total)
	}
}

func TestWordRepository_Update(t *testing.T) {
	// db, cleanup := setupWordTestDB(t)
	// defer cleanup()

	repo := service.NewMockWordRepository()
	ctx := context.Background()

	// Create a test word
	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts:   map[string]any{"verb_type": "ru-verb"},
	}
	err := repo.Create(ctx, word)
	if err != nil {
		t.Fatalf("error creating test word: %v", err)
	}

	// Update the word
	word.English = "to consume"
	err = repo.Update(ctx, word)
	if err != nil {
		t.Errorf("error updating word: %v", err)
	}

	// Verify the update
	updatedWord, err := repo.GetByID(ctx, word.ID)
	if err != nil {
		t.Errorf("error getting updated word: %v", err)
	}

	if updatedWord.English != "to consume" {
		t.Errorf("expected updated English to be 'to consume', got %s", updatedWord.English)
	}
}

func TestWordRepository_Delete(t *testing.T) {

	repo := service.NewMockWordRepository()
	ctx := context.Background()

	// Create a test word
	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts:   map[string]any{"verb_type": "ru-verb"},
	}
	err := repo.Create(ctx, word)
	if err != nil {
		t.Fatalf("error creating test word: %v", err)
	}

	// Delete the word
	err = repo.Delete(ctx, word.ID)
	if err != nil {
		t.Errorf("error deleting word: %v", err)
	}

	// Verify deletion
	deletedWord, err := repo.GetByID(ctx, word.ID)
	if err != nil {
		t.Errorf("error getting deleted word: %v", err)
	}

	if deletedWord != nil {
		t.Fatal("expected word to be deleted, got non-nil")
	}
} 