package service

import (
	"context"
	"fmt"
	"net/url"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository"
)

type StudyActivityService struct {
	activityRepo repository.StudyActivityRepository
	sessionRepo  repository.StudySessionRepository
}

func NewStudyActivityService(activityRepo repository.StudyActivityRepository, sessionRepo repository.StudySessionRepository) *StudyActivityService {
	return &StudyActivityService{
		activityRepo: activityRepo,
		sessionRepo:  sessionRepo,
	}
}

func (s *StudyActivityService) ListActivities(ctx context.Context) ([]*models.StudyActivity, error) {
	activities, err := s.activityRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing study activities: %v", err)
	}
	return activities, nil
}

func (s *StudyActivityService) GetActivity(ctx context.Context, id int64) (*models.StudyActivity, error) {
	activity, err := s.activityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting study activity: %v", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found")
	}
	return activity, nil
}

func (s *StudyActivityService) CreateActivity(ctx context.Context, activity *models.StudyActivity) error {
	if err := validateActivity(activity); err != nil {
		return err
	}

	if err := s.activityRepo.Create(ctx, activity); err != nil {
		return fmt.Errorf("error creating study activity: %v", err)
	}

	return nil
}

func (s *StudyActivityService) UpdateActivity(ctx context.Context, activity *models.StudyActivity) error {
	if err := validateActivity(activity); err != nil {
		return err
	}

	if err := s.activityRepo.Update(ctx, activity); err != nil {
		return fmt.Errorf("error updating study activity: %v", err)
	}

	return nil
}

func (s *StudyActivityService) DeleteActivity(ctx context.Context, id int64) error {
	if err := s.activityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting study activity: %v", err)
	}
	return nil
}

func validateActivity(activity *models.StudyActivity) error {
	if activity.Name == "" {
		return fmt.Errorf("activity name is required")
	}

	if activity.URL == "" {
		return fmt.Errorf("activity URL is required")
	}

	// Validate URL format
	_, err := url.ParseRequestURI(activity.URL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	return nil
}

func (s *StudyActivityService) ListActivitySessions(ctx context.Context, activityID int64, page, pageSize int) ([]*models.StudySession, error) {
	return s.sessionRepo.ListByActivity(ctx, activityID, page, pageSize)
} 