package implementations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository/sqlite"
)

type StudySessionRepository struct {
	db *sqlite.Database
}

func NewStudySessionRepository(db *sqlite.Database) *StudySessionRepository {
	return &StudySessionRepository{db: db}
}

func (r *StudySessionRepository) Create(ctx context.Context, session *models.StudySession) error {
	query := `
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		session.GroupID,
		session.StudyActivityID,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating study session: %v", err)
	}

	return nil
}

func (r *StudySessionRepository) GetByID(ctx context.Context, id int64) (*models.StudySession, error) {
	query := `
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions 
		WHERE id = ?`

	session := &models.StudySession{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}

	return session, nil
}

func (r *StudySessionRepository) ListByGroup(ctx context.Context, groupID int64, page, pageSize int) ([]*models.StudySession, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE group_id = ?`

	err := r.db.QueryRowContext(ctx, countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting study sessions: %v", err)
	}

	// Main query
	query := `
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, groupID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []*models.StudySession
	for rows.Next() {
		session := &models.StudySession{}
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.StudyActivityID,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning study session: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating study sessions: %v", err)
	}

	return sessions, total, nil
}

func (r *StudySessionRepository) AddReview(ctx context.Context, review *models.WordReviewItem) error {
	query := `
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		review.WordID,
		review.StudySessionID,
		review.Correct,
	).Scan(&review.ID, &review.CreatedAt)

	if err != nil {
		return fmt.Errorf("error adding word review: %v", err)
	}

	return nil
}

func (r *StudySessionRepository) GetSessionStats(ctx context.Context, sessionID int64) (*models.StudySessionStats, error) {
	query := `
		WITH session_reviews AS (
			SELECT 
				MIN(created_at) as start_time,
				MAX(created_at) as end_time,
				COUNT(*) as total_reviews,
				SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_reviews
			FROM word_review_items
			WHERE study_session_id = ?
		)
		SELECT 
			total_reviews,
			correct_reviews,
			CAST(
				(JULIANDAY(end_time) - JULIANDAY(start_time)) * 24 * 60 
				AS INTEGER
			) as duration_minutes
		FROM session_reviews`

	stats := &models.StudySessionStats{}
	var correctReviews int

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&stats.TotalReviews,
		&correctReviews,
		&stats.DurationMinutes,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting session stats: %v", err)
	}

	stats.CorrectReviews = correctReviews
	if stats.TotalReviews > 0 {
		stats.Accuracy = float64(correctReviews) / float64(stats.TotalReviews) * 100
	}

	return stats, nil
}

func (r *StudySessionRepository) ListReviews(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error) {
	query := `
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error listing reviews: %v", err)
	}
	defer rows.Close()

	var reviews []*models.WordReviewItem
	for rows.Next() {
		review := &models.WordReviewItem{}
		err := rows.Scan(
			&review.ID,
			&review.WordID,
			&review.StudySessionID,
			&review.Correct,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning review: %v", err)
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %v", err)
	}

	return reviews, nil
}

// FullReset deletes all study session related data
func (r *StudySessionRepository) FullReset(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tables := []string{
		"word_reviews",
		"session_words",
		"study_sessions",
		"study_activities",
	}
	
	for _, table := range tables {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to clear table %s: %v", table, err)
		}
	}

	return tx.Commit()
}

// GetLastSession retrieves the most recent study session with stats
func (r *StudySessionRepository) GetLastSession(ctx context.Context) (*models.StudySessionWithStats, error) {
	query := `
		SELECT s.id, s.created_at, 
			   a.id, a.name, a.url,
			   COUNT(r.id) as total_reviews,
			   SUM(CASE WHEN r.correct THEN 1 ELSE 0 END) as correct_reviews
		FROM study_sessions s
		JOIN study_activities a ON s.study_activity_id = a.id
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		GROUP BY s.id
		ORDER BY s.created_at DESC
		LIMIT 1`

	var s models.StudySessionWithStats
	var totalReviews, correctReviews int
	err := r.db.QueryRowContext(ctx, query).Scan(
		&s.ID, &s.CreatedAt,
		&s.StudyActivity.ID, &s.StudyActivity.Name, &s.StudyActivity.URL,
		&totalReviews, &correctReviews,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting last session: %v", err)
	}

	s.Stats.TotalReviews = totalReviews
	s.Stats.CorrectReviews = correctReviews
	if totalReviews > 0 {
		s.Stats.Accuracy = float64(correctReviews) / float64(totalReviews) * 100
	}

	return &s, nil
}

// GetQuickStats retrieves quick overview statistics
func (r *StudySessionRepository) GetQuickStats(ctx context.Context) (*models.QuickStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT s.id) as total_sessions,
			COUNT(r.id) as total_reviews,
			SUM(CASE WHEN r.correct THEN 1 ELSE 0 END) as correct_reviews
		FROM study_sessions s
		LEFT JOIN word_review_items r ON s.id = r.study_session_id`

	stats := &models.QuickStats{}
	var correctReviews int
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalSessions,
		&stats.TotalReviews,
		&correctReviews,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting quick stats: %v", err)
	}

	if stats.TotalReviews > 0 {
		stats.Accuracy = float64(correctReviews) / float64(stats.TotalReviews) * 100
	}

	return stats, nil
}

// GetSessionWords retrieves words associated with a specific study session
func (r *StudySessionRepository) GetSessionWords(ctx context.Context, sessionID int64) ([]*models.WordWithStats, error) {
	query := `
		SELECT w.id, w.kanji, w.romaji, w.english, w.parts,
			   COUNT(wr.id) as review_count, 
			   SUM(CASE WHEN wr.correct THEN 1 ELSE 0 END) as correct_count
		FROM words w
		JOIN session_words sw ON w.id = sw.word_id
		WHERE sw.session_id = ?
		GROUP BY w.id`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if err := json.Unmarshal(partsJSON, &w.Parts); err != nil {
			return nil, err
		}
		words = append(words, &w)
	}

	return words, nil
}

// GetStudyProgress retrieves study progress for the last 'days' days
func (r *StudySessionRepository) GetStudyProgress(ctx context.Context, days int) (*models.StudyProgress, error) {
	query := `
		SELECT 
			COUNT(DISTINCT s.id) as total_sessions,
			COUNT(r.id) as total_reviews,
			SUM(CASE WHEN r.correct THEN 1 ELSE 0 END) as correct_reviews
		FROM study_sessions s
		LEFT JOIN word_review_items r ON s.id = r.study_session_id
		WHERE s.created_at >= DATE('now', '-' || ? || ' days')`

	progress := &models.StudyProgress{}
	err := r.db.QueryRowContext(ctx, query, days).Scan(
		&progress.TotalSessions,
		&progress.TotalReviews,
		&progress.CorrectReviews,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting study progress: %v", err)
	}

	return progress, nil
}

// ListByActivity retrieves study sessions for a specific activity with pagination
func (r *StudySessionRepository) ListByActivity(ctx context.Context, activityID int64, page, pageSize int) ([]*models.StudySession, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Main query
	query := `
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, activityID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []*models.StudySession
	for rows.Next() {
		session := &models.StudySession{}
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.StudyActivityID,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating study sessions: %v", err)
	}

	return sessions, nil
}

// LoadSeedData reads and executes SQL files from the seeds directory
func (r *StudySessionRepository) LoadSeedData(ctx context.Context, seedsDir string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// List of seed files in order of execution
	seedFiles := []string{
		"study_sessions.sql", // Add your specific seed file for study sessions
	}

	for _, filename := range seedFiles {
		filepath := filepath.Join(seedsDir, filename)
		content, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %v", filename, err)
		}

		// Split content into individual SQL statements
		statements := strings.Split(string(content), ";")
		
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := tx.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("failed to execute seed query from %s: %v", filename, err)
			}
		}
	}

	return tx.Commit()
} 