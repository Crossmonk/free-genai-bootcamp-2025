package service

import (
	"context"
	"fmt"
	"math"

	"backend-go/internal/domain/models"
	"backend-go/internal/repository"
	"backend-go/internal/responses"
)

type GroupService struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(groupRepo repository.GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

type ListGroupsParams struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

type ListGroupsResult struct {
	Groups      []*models.Group
	TotalItems  int
	CurrentPage int
	TotalPages  int
}

func (s *GroupService) ListGroups(ctx context.Context, params ListGroupsParams) (*ListGroupsResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	groups, total, err := s.groupRepo.List(ctx, params.Page, params.PageSize, params.SortBy, params.Order)
	if err != nil {
		return nil, fmt.Errorf("error listing groups: %v", err)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &ListGroupsResult{
		Groups:      groups,
		TotalItems:  total,
		CurrentPage: params.Page,
		TotalPages:  totalPages,
	}, nil
}

func (s *GroupService) GetGroup(ctx context.Context, id int64) (*models.Group, error) {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found")
	}
	return group, nil
}

func (s *GroupService) GetGroupWithStats(ctx context.Context, id int64) (*models.GroupWithStats, error) {
	group, err := s.GetGroup(ctx, id)
	if err != nil {
		return nil, err
	}

	stats, err := s.groupRepo.GetStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting group stats: %v", err)
	}

	return &models.GroupWithStats{
		Group: *group,
		Stats: *stats,
	}, nil
}

func (s *GroupService) CreateGroup(ctx context.Context, group *models.Group) error {
	if err := validateGroup(group); err != nil {
		return err
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	return nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, group *models.Group) error {
	if err := validateGroup(group); err != nil {
		return err
	}

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return fmt.Errorf("error updating group: %v", err)
	}

	return nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, id int64) error {
	if err := s.groupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting group: %v", err)
	}
	return nil
}

type ListGroupWordsParams struct {
	GroupID  int64
	Page     int
	PageSize int
}

type ListGroupWordsResult struct {
	Words       []*models.WordWithStats
	TotalItems  int
	CurrentPage int
	TotalPages  int
}

func (s *GroupService) ListGroupWords(ctx context.Context, params ListGroupWordsParams) (*ListGroupWordsResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	words, total, err := s.groupRepo.ListWords(ctx, params.GroupID, params.Page, params.PageSize)
	if err != nil {
		return nil, fmt.Errorf("error listing group words: %v", err)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &ListGroupWordsResult{
		Words:       words,
		TotalItems:  total,
		CurrentPage: params.Page,
		TotalPages:  totalPages,
	}, nil
}

func (s *GroupService) AddWordToGroup(ctx context.Context, groupID, wordID int64) error {
	if err := s.groupRepo.AddWord(ctx, groupID, wordID); err != nil {
		return fmt.Errorf("error adding word to group: %v", err)
	}
	return nil
}

func (s *GroupService) RemoveWordFromGroup(ctx context.Context, groupID, wordID int64) error {
	if err := s.groupRepo.RemoveWord(ctx, groupID, wordID); err != nil {
		return fmt.Errorf("error removing word from group: %v", err)
	}
	return nil
}

func validateGroup(group *models.Group) error {
	if group.Name == "" {
		return fmt.Errorf("group name is required")
	}
	return nil
}

type GroupWordsResponse struct {
	GroupName string
	Words     []models.WordWithStats
}

// GetGroupWords retrieves all words belonging to a specific group with their review statistics
func (s *GroupService) GetGroupWords(ctx context.Context, groupID int64, page int, sortBy, order string) (*GroupWordsResponse, *responses.Pagination, error) {
	// Validate group exists
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	// Get paginated words with stats
	words, total, err := s.groupRepo.GetGroupWords(ctx, groupID, page, sortBy, order)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination
	itemsPerPage := 10 // Could be made configurable
	pagination := &responses.Pagination{
		CurrentPage:   page,
		ItemsPerPage:  itemsPerPage,
		TotalItems:    total,
		TotalPages:    int(math.Ceil(float64(total) / float64(itemsPerPage))),
	}

	// Convert []*WordWithStats to []WordWithStats
	wordsList := make([]models.WordWithStats, len(words))
	for i, w := range words {
		wordsList[i] = *w
	}

	return &GroupWordsResponse{
		GroupName: group.Name,
		Words:     wordsList,
	}, pagination, nil
}

type GroupStudySessionsResponse struct {
	GroupName     string
	StudySessions []models.StudySessionWithStats
	CurrentPage   int
	TotalPages    int
	TotalItems    int
}

// GetGroupStudySessions retrieves study sessions for a specific group with performance metrics
func (s *GroupService) GetGroupStudySessions(ctx context.Context, groupID int64, page, pageSize int) (*GroupStudySessionsResponse, error) {
	// Validate group exists
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("error getting group: %v", err)
	}

	// Get paginated study sessions
	sessions, total, err := s.groupRepo.ListStudySessions(ctx, groupID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("error listing study sessions: %v", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &GroupStudySessionsResponse{
		GroupName:     group.Name,
		StudySessions: sessions,
		CurrentPage:   page,
		TotalPages:    totalPages,
		TotalItems:    total,
	}, nil
} 