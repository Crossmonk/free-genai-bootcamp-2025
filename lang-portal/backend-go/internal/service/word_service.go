package service

import (
	"context"
	"fmt"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository"
)

type WordService struct {
	wordRepo repository.WordRepository
}

func NewWordService(wordRepo repository.WordRepository) *WordService {
	return &WordService{
		wordRepo: wordRepo,
	}
}

type WordServiceLite struct {
	wordRepo repository.WordRepository
}

func NewWordServiceLite(wordRepo repository.WordRepository) *WordServiceLite {
	return &WordServiceLite{wordRepo: wordRepo}
}

type ListWordsParams struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

type ListWordsResult struct {
	Words       []*models.WordWithStats
	TotalItems  int
	CurrentPage int
	TotalPages  int
}

func (s *WordService) ListWords(ctx context.Context, params ListWordsParams) (*ListWordsResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	words, total, err := s.wordRepo.List(ctx, params.Page, params.PageSize, params.SortBy, params.Order)
	if err != nil {
		return nil, fmt.Errorf("error listing words: %v", err)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &ListWordsResult{
		Words:       words,
		TotalItems:  total,
		CurrentPage: params.Page,
		TotalPages:  totalPages,
	}, nil
}

func (s *WordService) GetWord(ctx context.Context, id int64) (*models.Word, error) {
	word, err := s.wordRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting word: %v", err)
	}
	if word == nil {
		return nil, fmt.Errorf("word not found")
	}
	return word, nil
}

func (s *WordService) GetWordWithStats(ctx context.Context, id int64) (*models.WordWithStats, error) {
	word, err := s.GetWord(ctx, id)
	if err != nil {
		return nil, err
	}

	stats, err := s.wordRepo.GetStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting word stats: %v", err)
	}

	return &models.WordWithStats{
		ID:      word.ID,
		Kanji:   word.Kanji,
		Romaji:  word.Romaji,
		English: word.English,
		Parts:   word.Parts,
		Stats:   *stats,
	}, nil
}

func (s *WordService) CreateWord(ctx context.Context, word *models.Word) error {
	if err := validateWord(word); err != nil {
		return err
	}

	if err := s.wordRepo.Create(ctx, word); err != nil {
		return fmt.Errorf("error creating word: %v", err)
	}

	return nil
}

func (s *WordService) UpdateWord(ctx context.Context, word *models.Word) error {
	if err := validateWord(word); err != nil {
		return err
	}

	if err := s.wordRepo.Update(ctx, word); err != nil {
		return fmt.Errorf("error updating word: %v", err)
	}

	return nil
}

func (s *WordService) DeleteWord(ctx context.Context, id int64) error {
	if err := s.wordRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting word: %v", err)
	}
	return nil
}

func validateWord(word *models.Word) error {
	if word.Kanji == "" {
		return fmt.Errorf("kanji is required")
	}
	if word.Romaji == "" {
		return fmt.Errorf("romaji is required")
	}
	if word.English == "" {
		return fmt.Errorf("english translation is required")
	}
	if word.Parts == nil {
		return fmt.Errorf("parts is required")
	}
	return nil
} 