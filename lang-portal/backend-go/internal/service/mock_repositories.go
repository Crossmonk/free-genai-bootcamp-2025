package service

import (
	"context"
	"fmt"

	"backend-go/internal/domain/models"
)

type mockWordRepository struct {
	words map[int64]*models.Word
	stats map[int64]*models.WordStats
}

func NewMockWordRepository() *mockWordRepository {
	return &mockWordRepository{
		words: make(map[int64]*models.Word),
		stats: make(map[int64]*models.WordStats),
	}
}

func (m *mockWordRepository) Create(ctx context.Context, word *models.Word) error {
	word.ID = int64(len(m.words) + 1)
	m.words[word.ID] = word
	m.stats[word.ID] = &models.WordStats{}
	return nil
}

func (m *mockWordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	word, exists := m.words[id]
	if !exists {
		return nil, nil
	}
	return word, nil
}

func (m *mockWordRepository) List(ctx context.Context, page, pageSize int, sortBy, order string) ([]*models.WordWithStats, int, error) {
	var words []*models.WordWithStats
	for _, word := range m.words {
		stats := m.stats[word.ID]
		if stats == nil {
			stats = &models.WordStats{}
		}
		words = append(words, &models.WordWithStats{
			ID:      word.ID,
			Kanji:   word.Kanji,
			Romaji:  word.Romaji,
			English: word.English,
			Parts:   word.Parts,
			Stats:   *stats,
		})
	}
	return words, len(words), nil
}

func (m *mockWordRepository) Update(ctx context.Context, word *models.Word) error {
	if _, exists := m.words[word.ID]; !exists {
		return fmt.Errorf("word not found")
	}
	m.words[word.ID] = word
	return nil
}

func (m *mockWordRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.words[id]; !exists {
		return fmt.Errorf("word not found")
	}
	delete(m.words, id)
	delete(m.stats, id)
	return nil
}

func (m *mockWordRepository) GetStats(ctx context.Context, wordID int64) (*models.WordStats, error) {
	stats, exists := m.stats[wordID]
	if !exists {
		return &models.WordStats{}, nil
	}
	return stats, nil
}

func (m *mockWordRepository) UpdateStats(ctx context.Context, wordID int64, correct bool) error {
	stats, exists := m.stats[wordID]
	if !exists {
		stats = &models.WordStats{}
		m.stats[wordID] = stats
	}

	if correct {
		stats.CorrectCount++
	} else {
		stats.WrongCount++
	}

	total := stats.CorrectCount + stats.WrongCount
	if total > 0 {
		stats.Accuracy = float64(stats.CorrectCount) / float64(total) * 100
	}

	return nil
}

type mockGroupRepository struct {
	groups     map[int64]*models.Group
	stats      map[int64]*models.GroupStats
	wordGroups map[int64]map[int64]bool
}

func NewMockGroupRepository() *mockGroupRepository {
	return &mockGroupRepository{
		groups:     make(map[int64]*models.Group),
		stats:      make(map[int64]*models.GroupStats),
		wordGroups: make(map[int64]map[int64]bool),
	}
}

func (m *mockGroupRepository) AddWord(ctx context.Context, groupID, wordID int64) error {
	if _, exists := m.groups[groupID]; !exists {
		return fmt.Errorf("group not found")
	}
	if m.wordGroups[groupID] == nil {
		m.wordGroups[groupID] = make(map[int64]bool)
	}
	m.wordGroups[groupID][wordID] = true
	return nil
}

func (m *mockGroupRepository) RemoveWord(ctx context.Context, groupID, wordID int64) error {
	if _, exists := m.groups[groupID]; !exists {
		return fmt.Errorf("group not found")
	}
	delete(m.wordGroups[groupID], wordID)
	return nil
}

func (m *mockGroupRepository) ListWords(ctx context.Context, groupID int64) ([]int64, error) {
	if _, exists := m.groups[groupID]; !exists {
		return nil, fmt.Errorf("group not found")
	}
	var wordIDs []int64
	for wordID := range m.wordGroups[groupID] {
		wordIDs = append(wordIDs, wordID)
	}
	return wordIDs, nil
}

func (m *mockGroupRepository) Create(ctx context.Context, group *models.Group) error {
	group.ID = int64(len(m.groups) + 1)
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepository) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	group, exists := m.groups[id]
	if !exists {
		return nil, nil
	}
	return group, nil
}

func (m *mockGroupRepository) List(ctx context.Context, page, pageSize int, sortBy, order string) ([]*models.Group, int, error) {
	var groups []*models.Group
	for _, group := range m.groups {
		groups = append(groups, group)
	}
	return groups, len(groups), nil
}

func (m *mockGroupRepository) Update(ctx context.Context, group *models.Group) error {
	if _, exists := m.groups[group.ID]; !exists {
		return fmt.Errorf("group not found")
	}
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.groups[id]; !exists {
		return fmt.Errorf("group not found")
	}
	delete(m.groups, id)
	return nil
}

func (m *mockGroupRepository) GetStats(ctx context.Context, groupID int64) (*models.GroupStats, error) {
	stats, exists := m.stats[groupID]
	if !exists {
		return &models.GroupStats{}, nil
	}
	return stats, nil
}

func NewMockStudyActivityRepository() *mockStudyActivityRepository {
	return &mockStudyActivityRepository{
		activities: make(map[int64]*models.StudyActivity),
	}
}

type mockStudyActivityRepository struct {
	activities map[int64]*models.StudyActivity
}

func (m *mockStudyActivityRepository) Create(ctx context.Context, activity *models.StudyActivity) error {
	activity.ID = int64(len(m.activities) + 1)
	m.activities[activity.ID] = activity
	return nil
}

func (m *mockStudyActivityRepository) GetByID(ctx context.Context, id int64) (*models.StudyActivity, error) {
	activity, exists := m.activities[id]
	if !exists {
		return nil, nil
	}
	return activity, nil
}

func (m *mockStudyActivityRepository) List(ctx context.Context) ([]*models.StudyActivity, error) {
	var activities []*models.StudyActivity
	for _, activity := range m.activities {
		activities = append(activities, activity)
	}
	return activities, nil
}

func (m *mockStudyActivityRepository) Update(ctx context.Context, activity *models.StudyActivity) error {
	if _, exists := m.activities[activity.ID]; !exists {
		return fmt.Errorf("activity not found")
	}
	m.activities[activity.ID] = activity
	return nil
}

func (m *mockStudyActivityRepository) Delete(ctx context.Context, id int64) error {
	if _, exists := m.activities[id]; !exists {
		return fmt.Errorf("activity not found")
	}
	delete(m.activities, id)
	return nil
}

type mockStudySessionRepository struct {
	sessions map[int64]*models.StudySession
	reviews  map[int64][]*models.WordReviewItem
}

func NewMockStudySessionRepository() *mockStudySessionRepository {
	return &mockStudySessionRepository{
		sessions: make(map[int64]*models.StudySession),
		reviews:  make(map[int64][]*models.WordReviewItem),
	}
}

func (m *mockStudySessionRepository) Create(ctx context.Context, session *models.StudySession) error {
	session.ID = int64(len(m.sessions) + 1)
	m.sessions[session.ID] = session
	return nil
}

func (m *mockStudySessionRepository) GetByID(ctx context.Context, id int64) (*models.StudySession, error) {
	session, exists := m.sessions[id]
	if !exists {
		return nil, nil
	}
	return session, nil
}

func (m *mockStudySessionRepository) ListByGroup(ctx context.Context, groupID int64, page, pageSize int) ([]*models.StudySession, int, error) {
	var sessions []*models.StudySession
	for _, session := range m.sessions {
		if session.GroupID == groupID {
			sessions = append(sessions, session)
		}
	}
	return sessions, len(sessions), nil
}

func (m *mockStudySessionRepository) AddReview(ctx context.Context, review *models.WordReviewItem) error {
	review.ID = int64(len(m.reviews[review.StudySessionID]) + 1)
	m.reviews[review.StudySessionID] = append(m.reviews[review.StudySessionID], review)
	return nil
}

func (m *mockStudySessionRepository) ListReviews(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error) {
	reviews := m.reviews[sessionID]
	if reviews == nil {
		reviews = []*models.WordReviewItem{}
	}
	return reviews, nil
}

func (m *mockStudySessionRepository) GetSessionStats(ctx context.Context, sessionID int64) (*models.StudySessionStats, error) {
	reviews := m.reviews[sessionID]
	stats := &models.StudySessionStats{
		TotalReviews: len(reviews),
	}
	for _, review := range reviews {
		if review.Correct {
			stats.CorrectReviews++
		}
	}
	if stats.TotalReviews > 0 {
		stats.Accuracy = float64(stats.CorrectReviews) / float64(stats.TotalReviews) * 100
	}
	return stats, nil
}

func (m *mockStudySessionRepository) FullReset(ctx context.Context) error {
	m.sessions = make(map[int64]*models.StudySession)
	m.reviews = make(map[int64][]*models.WordReviewItem)
	return nil
}

func (m *mockStudySessionRepository) GetLastSession(ctx context.Context) (*models.StudySessionWithStats, error) {
	if len(m.sessions) == 0 {
		return nil, nil // No sessions available
	}

	// Find the last session (assuming sessions are added in order)
	var lastSession *models.StudySession
	for _, session := range m.sessions {
		if lastSession == nil || session.CreatedAt.After(lastSession.CreatedAt) {
			lastSession = session
		}
	}

	// Create a mock StudySessionWithStats
	stats := &models.StudySessionStats{
		TotalReviews: len(m.reviews[lastSession.ID]),
	}
	for _, review := range m.reviews[lastSession.ID] {
		if review.Correct {
			stats.CorrectReviews++
		}
	}
	if stats.TotalReviews > 0 {
		stats.Accuracy = float64(stats.CorrectReviews) / float64(stats.TotalReviews) * 100
	}

	return &models.StudySessionWithStats{
		ID:     lastSession.ID,
		Stats:  *stats,
		// Populate other fields as necessary
	}, nil
}

func (m *mockStudySessionRepository) GetQuickStats(ctx context.Context) (*models.QuickStats, error) {
	stats := &models.QuickStats{
		TotalSessions: len(m.sessions),
		TotalReviews:  0,
	}

	for _, reviews := range m.reviews {
		stats.TotalReviews += len(reviews)
	}

	if stats.TotalReviews > 0 {
		stats.Accuracy = float64(stats.TotalReviews) / float64(stats.TotalReviews) * 100
	}

	return stats, nil
}

func (m *mockStudySessionRepository) GetSessionWords(ctx context.Context, sessionID int64) ([]*models.WordWithStats, error) {
	// For simplicity, return an empty slice or mock data
	// You can modify this to return actual mock data if needed
	return []*models.WordWithStats{}, nil
}

func (m *mockStudySessionRepository) GetStudyProgress(ctx context.Context, days int) (*models.StudyProgress, error) {
	// For simplicity, return mock data
	// You can modify this to return actual mock data if needed
	return &models.StudyProgress{
		TotalSessions: len(m.sessions),
		TotalReviews:  0, // Adjust as needed
		CorrectReviews: 0, // Adjust as needed
	}, nil
}

func (m *mockStudySessionRepository) ListByActivity(ctx context.Context, activityID int64, page, pageSize int) ([]*models.StudySession, error) {
	var sessions []*models.StudySession
	for _, session := range m.sessions {
		if session.StudyActivityID == activityID {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (m *mockStudySessionRepository) LoadSeedData(ctx context.Context, seedsDir string) error {
	// For simplicity, you can just return nil or implement logic to load mock data
	// Here, we will just simulate loading some mock data
	m.sessions[1] = &models.StudySession{ID: 1, StudyActivityID: 1, GroupID: 1}
	m.sessions[2] = &models.StudySession{ID: 2, StudyActivityID: 1, GroupID: 1}
	return nil
}

