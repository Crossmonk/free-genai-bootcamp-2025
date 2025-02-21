package service

import (
	"context"
	"testing"

	"backend-go/internal/domain/models"
)


func newMockWordRepository() *mockWordRepository {
	return &mockWordRepository{
		words: make(map[int64]*models.Word),
		stats: make(map[int64]*models.WordStats),
	}
}


func TestWordService_CreateWord(t *testing.T) {
	repo := newMockWordRepository()
	service := NewWordService(repo)
	ctx := context.Background()

	tests := []struct {
		name    string
		word    *models.Word
		wantErr bool
	}{
		{
			name: "valid word",
			word: &models.Word{
				Kanji:   "食べる",
				Romaji:  "taberu",
				English: "to eat",
				Parts: map[string]any{
					"verb_type": "ru-verb",
					"topic":     "food",
				},
			},
			wantErr: false,
		},
		{
			name: "missing kanji",
			word: &models.Word{
				Romaji:  "taberu",
				English: "to eat",
				Parts:   map[string]any{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateWord(ctx, tt.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.word.ID == 0 {
				t.Error("CreateWord() did not set ID for valid word")
			}
		})
	}
}

func TestWordService_GetWordWithStats(t *testing.T) {
	repo := newMockWordRepository()
	service := NewWordService(repo)
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
	err := service.CreateWord(ctx, word)
	if err != nil {
		t.Fatalf("Failed to create test word: %v", err)
	}

	// Set some stats
	repo.stats[word.ID] = &models.WordStats{
		CorrectCount: 5,
		WrongCount:   2,
		Accuracy:     71.43,
	}

	// Test getting the word with stats
	wordWithStats, err := service.GetWordWithStats(ctx, word.ID)
	if err != nil {
		t.Errorf("GetWordWithStats() error = %v", err)
		return
	}

	if wordWithStats.ID != word.ID {
		t.Errorf("GetWordWithStats() got ID = %v, want %v", wordWithStats.ID, word.ID)
	}

	if wordWithStats.Stats.CorrectCount != 5 {
		t.Errorf("GetWordWithStats() got CorrectCount = %v, want %v", wordWithStats.Stats.CorrectCount, 5)
	}
} 