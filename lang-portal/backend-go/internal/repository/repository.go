package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"backend-go/internal/domain/models"
)

type WordRepository interface {
	Create(ctx context.Context, word *models.Word) error
	GetByID(ctx context.Context, id int64) (*models.Word, error)
	List(ctx context.Context, page, pageSize int, sortBy, order string) ([]*models.WordWithStats, int, error)
	Update(ctx context.Context, word *models.Word) error
	Delete(ctx context.Context, id int64) error
	GetStats(ctx context.Context, wordID int64) (*models.WordStats, error)
}

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetByID(ctx context.Context, id int64) (*models.Group, error)
	List(ctx context.Context, page, pageSize int, sortBy, order string) ([]*models.Group, int, error)
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id int64) error
	GetStats(ctx context.Context, groupID int64) (*models.GroupStats, error)
	AddWord(ctx context.Context, groupID, wordID int64) error
	RemoveWord(ctx context.Context, groupID, wordID int64) error
	ListWords(ctx context.Context, groupID int64, page, pageSize int) ([]*models.WordWithStats, int, error)
	GetGroupWords(ctx context.Context, groupID int64, page int, sortBy, order string) ([]*models.WordWithStats, int, error)
	ListStudySessions(ctx context.Context, groupID int64, page, pageSize int) ([]models.StudySessionWithStats, int, error)
}

type StudyActivityRepository interface {
	Create(ctx context.Context, activity *models.StudyActivity) error
	GetByID(ctx context.Context, id int64) (*models.StudyActivity, error)
	List(ctx context.Context) ([]*models.StudyActivity, error)
	Update(ctx context.Context, activity *models.StudyActivity) error
	Delete(ctx context.Context, id int64) error
}

type StudySessionRepository interface {
	Create(ctx context.Context, session *models.StudySession) error
	GetByID(ctx context.Context, id int64) (*models.StudySession, error)
	ListByGroup(ctx context.Context, groupID int64, page, pageSize int) ([]*models.StudySession, int, error)
	AddReview(ctx context.Context, review *models.WordReviewItem) error
	GetSessionStats(ctx context.Context, sessionID int64) (*models.StudySessionStats, error)
	ListReviews(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error)
	GetLastSession(ctx context.Context) (*models.StudySessionWithStats, error)
	GetStudyProgress(ctx context.Context, days int) (*models.StudyProgress, error)
	GetQuickStats(ctx context.Context) (*models.QuickStats, error)
	GetSessionWords(ctx context.Context, sessionID int64) ([]*models.WordWithStats, error)
	FullReset(ctx context.Context) error
	LoadSeedData(ctx context.Context, seedsDir string) error
	ListByActivity(ctx context.Context, activityID int64, page, pageSize int) ([]*models.StudySession, error)
}

type Repository struct {
	db *sql.DB
}

// FullReset deletes all study related data from the database
func (r *Repository) FullReset(ctx context.Context) error {
	// Using transaction to ensure all or nothing
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	// Order matters due to foreign key constraints
	tables := []string{
		"word_reviews",
		"session_words",
		"study_sessions",
		"study_activities",
	}
	
	for _, table := range tables {
		_, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to clear table %s: %v", table, err)
		}
	}
	
	return tx.Commit()
}

// LoadSeedData reads and executes SQL files from the seeds directory
func (r *Repository) LoadSeedData(ctx context.Context, seedsDir string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// List of seed files in order of execution
	seedFiles := []string{
		"activities.sql",
		"groups.sql",
		"words.sql",
		"group_words.sql",
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