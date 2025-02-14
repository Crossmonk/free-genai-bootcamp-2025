package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository/sqlite"
)

type StudyActivityRepository struct {
	db *sqlite.Database
}

func NewStudyActivityRepository(db *sqlite.Database) *StudyActivityRepository {
	return &StudyActivityRepository{db: db}
}

func (r *StudyActivityRepository) Create(ctx context.Context, activity *models.StudyActivity) error {
	query := `
		INSERT INTO study_activities (name, url)
		VALUES (?, ?)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		activity.Name,
		activity.URL,
	).Scan(&activity.ID)

	if err != nil {
		return fmt.Errorf("error creating study activity: %v", err)
	}

	return nil
}

func (r *StudyActivityRepository) GetByID(ctx context.Context, id int64) (*models.StudyActivity, error) {
	query := `SELECT id, name, url FROM study_activities WHERE id = ?`

	activity := &models.StudyActivity{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.Name,
		&activity.URL,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study activity: %v", err)
	}

	return activity, nil
}

func (r *StudyActivityRepository) List(ctx context.Context) ([]*models.StudyActivity, error) {
	query := `
		SELECT 
			sa.id, sa.name, sa.url,
			MAX(ss.created_at) as last_used
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON sa.id = ss.study_activity_id
		GROUP BY sa.id
		ORDER BY sa.name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing study activities: %v", err)
	}
	defer rows.Close()

	var activities []*models.StudyActivity
	for rows.Next() {
		activity := &models.StudyActivity{}
		var lastUsed sql.NullTime

		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.URL,
			&lastUsed,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study activity: %v", err)
		}

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating study activities: %v", err)
	}

	return activities, nil
}

func (r *StudyActivityRepository) Update(ctx context.Context, activity *models.StudyActivity) error {
	query := `
		UPDATE study_activities 
		SET name = ?, url = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		activity.Name,
		activity.URL,
		activity.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating study activity: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("study activity not found")
	}

	return nil
}

func (r *StudyActivityRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM study_activities WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting study activity: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("study activity not found")
	}

	return nil
} 