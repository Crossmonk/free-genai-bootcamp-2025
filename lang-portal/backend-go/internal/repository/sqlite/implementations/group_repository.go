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

type GroupRepository struct {
	db *sqlite.Database
}

func NewGroupRepository(db *sqlite.Database) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *models.Group) error {
	query := `
		INSERT INTO groups (name, words_count)
		VALUES (?, 0)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, group.Name).Scan(&group.ID)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	return nil
}

func (r *GroupRepository) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	query := `
		SELECT g.id, g.name, g.words_count, 
		       MAX(s.created_at) as last_studied_at
		FROM groups g
		LEFT JOIN study_sessions s ON g.id = s.group_id
		WHERE g.id = ?
		GROUP BY g.id`

	group := &models.Group{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.WordsCount,
		&group.LastStudiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting group: %v", err)
	}

	return group, nil
}

func (r *GroupRepository) List(ctx context.Context, page, pageSize int, sortBy, order string) ([]*models.Group, int, error) {
	// Validate and sanitize sort parameters
	allowedSortFields := map[string]string{
		"name":        "g.name",
		"words_count": "g.words_count",
	}

	dbSortField, ok := allowedSortFields[sortBy]
	if !ok {
		dbSortField = "g.name"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM groups`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting groups: %v", err)
	}

	// Main query
	query := `
		SELECT g.id, g.name, g.words_count, 
		       MAX(s.created_at) as last_studied_at
		FROM groups g
		LEFT JOIN study_sessions s ON g.id = s.group_id
		GROUP BY g.id
		ORDER BY ` + dbSortField + ` ` + order + `
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing groups: %v", err)
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group := &models.Group{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.WordsCount,
			&group.LastStudiedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning group: %v", err)
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating groups: %v", err)
	}

	return groups, total, nil
}

func (r *GroupRepository) Update(ctx context.Context, group *models.Group) error {
	query := `
		UPDATE groups 
		SET name = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, group.Name, group.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}

func (r *GroupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM groups WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}

func (r *GroupRepository) GetStats(ctx context.Context, groupID int64) (*models.GroupStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_reviews,
			SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_reviews
		FROM word_review_items wri
		JOIN study_sessions ss ON wri.study_session_id = ss.id
		WHERE ss.group_id = ?`

	stats := &models.GroupStats{}
	var correctReviews int
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(
		&stats.TotalReviews,
		&correctReviews,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting group stats: %v", err)
	}

	stats.CorrectReviews = correctReviews
	if stats.TotalReviews > 0 {
		stats.Accuracy = float64(correctReviews) / float64(stats.TotalReviews) * 100
	}

	return stats, nil
}

func (r *GroupRepository) AddWord(ctx context.Context, groupID, wordID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Add word to group
	query := `INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, query, wordID, groupID)
	if err != nil {
		return fmt.Errorf("error adding word to group: %v", err)
	}

	// Update words count
	updateQuery := `
		UPDATE groups 
		SET words_count = (
			SELECT COUNT(*) 
			FROM word_groups 
			WHERE group_id = ?
		)
		WHERE id = ?`

	_, err = tx.ExecContext(ctx, updateQuery, groupID, groupID)
	if err != nil {
		return fmt.Errorf("error updating words count: %v", err)
	}

	return tx.Commit()
}

func (r *GroupRepository) RemoveWord(ctx context.Context, groupID, wordID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Remove word from group
	query := `DELETE FROM word_groups WHERE word_id = ? AND group_id = ?`
	result, err := tx.ExecContext(ctx, query, wordID, groupID)
	if err != nil {
		return fmt.Errorf("error removing word from group: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found in group")
	}

	// Update words count
	updateQuery := `
		UPDATE groups 
		SET words_count = (
			SELECT COUNT(*) 
			FROM word_groups 
			WHERE group_id = ?
		)
		WHERE id = ?`

	_, err = tx.ExecContext(ctx, updateQuery, groupID, groupID)
	if err != nil {
		return fmt.Errorf("error updating words count: %v", err)
	}

	return tx.Commit()
}

func (r *GroupRepository) ListWords(ctx context.Context, groupID int64, page, pageSize int) ([]*models.WordWithStats, int, error) {
	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM word_groups 
		WHERE group_id = ?`
	
	err := r.db.QueryRowContext(ctx, countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting words in group: %v", err)
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Main query
	query := `
		SELECT 
			w.id, w.kanji, w.romaji, w.english, w.parts,
			COALESCE(correct_reviews.count, 0) as correct_count,
			COALESCE(wrong_reviews.count, 0) as wrong_count
		FROM words w
		JOIN word_groups wg ON w.id = wg.word_id
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
		WHERE wg.group_id = ?
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, groupID, pageSize, offset)
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

func (r *GroupRepository) ListStudySessions(ctx context.Context, groupID int64, page, pageSize int) ([]models.StudySessionWithStats, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM study_sessions WHERE group_id = ?`
	err := r.db.QueryRowContext(ctx, countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting study sessions: %v", err)
	}

	// Get study sessions with stats
	query := `
		SELECT 
			s.id, s.created_at,
			a.id, a.name, a.url,
			COUNT(r.id) as total_reviews,
			SUM(CASE WHEN r.correct THEN 1 ELSE 0 END) as correct_reviews
		FROM study_sessions s
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.group_id = ?
		GROUP BY s.id
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, groupID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []models.StudySessionWithStats
	for rows.Next() {
		var s models.StudySessionWithStats
		var totalReviews, correctReviews int
		err := rows.Scan(
			&s.ID, &s.CreatedAt,
			&s.StudyActivity.ID, &s.StudyActivity.Name, &s.StudyActivity.URL,
			&totalReviews, &correctReviews,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session: %v", err)
		}

		s.Stats.TotalReviews = totalReviews
		s.Stats.CorrectReviews = correctReviews
		if totalReviews > 0 {
			s.Stats.Accuracy = float64(correctReviews) / float64(totalReviews) * 100
		}
		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

// GetGroupWords retrieves paginated words with stats for a group
func (r *GroupRepository) GetGroupWords(ctx context.Context, groupID int64, page int, sortBy, order string) ([]*models.WordWithStats, int, error) {
	offset := (page - 1) * 10
	query := `
		SELECT w.id, w.kanji, w.romaji, w.english, w.parts,
			   COUNT(wr.id) as review_count, 
			   SUM(CASE WHEN wr.correct THEN 1 ELSE 0 END) as correct_count
		FROM words w
		JOIN group_words gw ON w.id = gw.word_id
		LEFT JOIN word_reviews wr ON w.id = wr.word_id
		WHERE gw.group_id = ?
		GROUP BY w.id
		ORDER BY ` + sortBy + ` ` + order + `
		LIMIT 10 OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, groupID, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []*models.WordWithStats
	for rows.Next() {
		var w models.WordWithStats
		var partsJSON []byte
		err := rows.Scan(
			&w.ID, &w.Kanji, &w.Romaji, &w.English, &partsJSON,
			&w.Stats.WrongCount, &w.Stats.CorrectCount, 
		)
		if err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(partsJSON, &w.Parts); err != nil {
			return nil, 0, err
		}
		words = append(words, &w)
	}

	var total int
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM group_words WHERE group_id = ?", groupID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return words, total, nil
} 