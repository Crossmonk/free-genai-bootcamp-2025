package service

import (
	"context"
	"fmt"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository"
)

type StudySessionService struct {
	sessionRepo repository.StudySessionRepository
	groupRepo   repository.GroupRepository
}

func NewStudySessionService(
	sessionRepo repository.StudySessionRepository,
	groupRepo repository.GroupRepository,
) *StudySessionService {
	return &StudySessionService{
		sessionRepo: sessionRepo,
		groupRepo:   groupRepo,
	}
}

type CreateSessionParams struct {
	GroupID         int64
	StudyActivityID int64
}

func (s *StudySessionService) CreateSession(ctx context.Context, params CreateSessionParams) (*models.StudySession, error) {
	// Verify group exists
	group, err := s.groupRepo.GetByID(ctx, params.GroupID)
	if err != nil {
		return nil, fmt.Errorf("error verifying group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}

	session := &models.StudySession{
		GroupID:         params.GroupID,
		StudyActivityID: params.StudyActivityID,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("error creating study session: %v", err)
	}

	return session, nil
}

func (s *StudySessionService) GetSession(ctx context.Context, id int64) (*models.StudySession, error) {
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found")
	}
	return session, nil
}

type ListSessionsParams struct {
	GroupID  int64
	Page     int
	PageSize int
}

type ListSessionsResult struct {
	Sessions    []*models.StudySession
	TotalItems  int
	CurrentPage int
	TotalPages  int
}

func (s *StudySessionService) ListSessions(ctx context.Context, params ListSessionsParams) (*ListSessionsResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	sessions, total, err := s.sessionRepo.ListByGroup(ctx, params.GroupID, params.Page, params.PageSize)
	if err != nil {
		return nil, fmt.Errorf("error listing study sessions: %v", err)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &ListSessionsResult{
		Sessions:    sessions,
		TotalItems:  total,
		CurrentPage: params.Page,
		TotalPages:  totalPages,
	}, nil
}

func (s *StudySessionService) AddReview(ctx context.Context, sessionID, wordID int64, correct bool) (*models.WordReviewItem, error) {
	// Verify session exists
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error verifying session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found")
	}

	review := &models.WordReviewItem{
		WordID:         wordID,
		StudySessionID: sessionID,
		Correct:        correct,
	}

	if err := s.sessionRepo.AddReview(ctx, review); err != nil {
		return nil, fmt.Errorf("error adding word review: %v", err)
	}

	return review, nil
}

func (s *StudySessionService) GetSessionStats(ctx context.Context, sessionID int64) (*models.StudySessionStats, error) {
	// Verify session exists
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error verifying session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found")
	}

	stats, err := s.sessionRepo.GetSessionStats(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting session stats: %v", err)
	}

	return stats, nil
}

func (s *StudySessionService) ListSessionReviews(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error) {
	// Verify session exists
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error verifying session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("study session not found")
	}

	reviews, err := s.sessionRepo.ListReviews(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error listing session reviews: %v", err)
	}

	return reviews, nil
}

func (s *StudySessionService) GetLastStudySession(ctx context.Context) (*models.StudySessionWithStats, error) {
	return s.sessionRepo.GetLastSession(ctx)
}

func (s *StudySessionService) GetStudyProgress(ctx context.Context, days int) (*models.StudyProgress, error) {
	return s.sessionRepo.GetStudyProgress(ctx, days)
}

func (s *StudySessionService) GetQuickStats(ctx context.Context) (*models.QuickStats, error) {
	return s.sessionRepo.GetQuickStats(ctx)
}

func (s *StudySessionService) GetSessionWords(ctx context.Context, sessionID int64) ([]*models.WordWithStats, error) {
	return s.sessionRepo.GetSessionWords(ctx, sessionID)
}

// LoadSeedData loads initial seed data from the seeds directory
func (s *StudySessionService) LoadSeedData(ctx context.Context) error {
	seedsDir := "seeds"  // relative to backend-go directory
	return s.sessionRepo.LoadSeedData(ctx, seedsDir)
}

// FullReset deletes all study session related data
func (s *StudySessionService) FullReset(ctx context.Context) error {
	return s.sessionRepo.FullReset(ctx)
}

func (s *StudySessionService) ListByActivity(ctx context.Context, activityID int64, page, pageSize int) ([]*models.StudySession, error) {
	sessions, err := s.sessionRepo.ListByActivity(ctx, activityID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error listing study sessions: %v", err)
	}

	return sessions, nil
} 