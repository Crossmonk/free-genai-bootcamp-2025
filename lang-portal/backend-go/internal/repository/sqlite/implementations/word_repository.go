package implementations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"backend-go/internal/domain/models"	
	"backend-go/internal/repository/sqlite"

)

type WordRepository struct {
	db *sqlite.Database
}

func NewWordRepository(db *sqlite.Database) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) Create(ctx context.Context, word *models.Word) error {
	parts, err := json.Marshal(word.Parts)
	if err != nil {
		return fmt.Errorf("error marshaling parts: %v", err)
	}

	query := `
		INSERT INTO words (kanji, romaji, english, parts)
		VALUES (?, ?, ?, ?)
		RETURNING id`

	err = r.db.QueryRowContext(ctx, query,
		word.Kanji,
		word.Romaji,
		word.English,
		parts,
	).Scan(&word.ID)

	if err != nil {
		return fmt.Errorf("error creating word: %v", err)
	}

	return nil
}

func (r *WordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	word := &models.Word{}
	var partsJSON []byte

	query := `SELECT id, kanji, romaji, english, parts FROM words WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID,
		&word.Kanji,
		&word.Romaji,
		&word.English,
		&partsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting word: %v", err)
	}

	if err := json.Unmarshal(partsJSON, &word.Parts); err != nil {
		return nil, fmt.Errorf("error unmarshaling parts: %v", err)
	}

	return word, nil
}

func (r *WordRepository) List(ctx context.Context, page, pageSize int, sortBy, order string) ([]*models.WordWithStats, int, error) {
	// Validate and sanitize sort parameters
	allowedSortFields := map[string]string{
		"kanji":         "w.kanji",
		"romaji":        "w.romaji",
		"english":       "w.english",
		"correct_count": "COALESCE(correct_reviews.count, 0)",
		"wrong_count":   "COALESCE(wrong_reviews.count, 0)",
	}

	dbSortField, ok := allowedSortFields[sortBy]
	if !ok {
		dbSortField = "w.kanji"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Query to get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM words`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting words: %v", err)
	}

	// Main query with stats
	query := `
		SELECT 
			w.id, w.kanji, w.romaji, w.english, w.parts,
			COALESCE(correct_reviews.count, 0) as correct_count,
			COALESCE(wrong_reviews.count, 0) as wrong_count
		FROM words w
		LEFT JOIN (
			SELECT word_id, COUNT(*) as count
			FROM word_review_items
			WHERE correct = true
			GROUP BY word_id
		) correct_reviews ON w.id = correct_reviews.word_id
		LEFT JOIN (
			SELECT word_id, COUNT(*) as count
			FROM word_review_items
			WHERE correct = false
			GROUP BY word_id
		) wrong_reviews ON w.id = wrong_reviews.word_id
		ORDER BY ` + dbSortField + ` ` + order + `
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing words: %v", err)
	}
	defer rows.Close()

	var words []*models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		var partsJSON []byte
		var correctCount, wrongCount int

		err := rows.Scan(
			&word.ID,
			&word.Kanji,
			&word.Romaji,
			&word.English,
			&partsJSON,
			&correctCount,
			&wrongCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word: %v", err)
		}

		if err := json.Unmarshal(partsJSON, &word.Parts); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling parts: %v", err)
		}

		// Calculate accuracy
		word.Stats = models.WordStats{
			CorrectCount: correctCount,
			WrongCount:   wrongCount,
		}
		
		totalAttempts := correctCount + wrongCount
		if totalAttempts > 0 {
			word.Stats.Accuracy = float64(correctCount) / float64(totalAttempts) * 100
		}

		words = append(words, &word)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating words: %v", err)
	}

	return words, total, nil
}

func (r *WordRepository) Update(ctx context.Context, word *models.Word) error {
	parts, err := json.Marshal(word.Parts)
	if err != nil {
		return fmt.Errorf("error marshaling parts: %v", err)
	}

	query := `
		UPDATE words 
		SET kanji = ?, romaji = ?, english = ?, parts = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		word.Kanji,
		word.Romaji,
		word.English,
		parts,
		word.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found")
	}

	return nil
}

func (r *WordRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM words WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found")
	}

	return nil
}

func (r *WordRepository) GetStats(ctx context.Context, wordID int64) (*models.WordStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN correct THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN NOT correct THEN 1 ELSE 0 END), 0) as wrong_count
		FROM word_review_items
		WHERE word_id = ?`

	stats := &models.WordStats{}
	err := r.db.QueryRowContext(ctx, query, wordID).Scan(
		&stats.CorrectCount,
		&stats.WrongCount,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting word stats: %v", err)
	}

	totalAttempts := stats.CorrectCount + stats.WrongCount
	if totalAttempts > 0 {
		stats.Accuracy = float64(stats.CorrectCount) / float64(totalAttempts) * 100
	}

	return stats, nil
} 