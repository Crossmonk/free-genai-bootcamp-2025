package handlers

import (
	"backend-go/internal/domain/models"
)

type ListWordsResponse struct {
	Data struct {
		Words      []*models.WordWithStats `json:"words"`
		Pagination struct {
			CurrentPage   int `json:"current_page"`
			TotalPages    int `json:"total_pages"`
			TotalItems    int `json:"total_items"`
			ItemsPerPage  int `json:"items_per_page"`
		} `json:"pagination"`
	} `json:"data"`
}

type WordResponse struct {
	Data models.WordWithStats `json:"data"`
}

type CreateWordRequest struct {
	Kanji   string         `json:"kanji" binding:"required"`
	Romaji  string         `json:"romaji" binding:"required"`
	English string         `json:"english" binding:"required"`
	Parts   map[string]any `json:"parts" binding:"required"`
}

type UpdateWordRequest struct {
	Kanji   string         `json:"kanji" binding:"required"`
	Romaji  string         `json:"romaji" binding:"required"`
	English string         `json:"english" binding:"required"`
	Parts   map[string]any `json:"parts" binding:"required"`
}

type ListGroupsResponse struct {
	Data struct {
		Groups     []*models.Group `json:"groups"`
		Pagination struct {
			CurrentPage   int `json:"current_page"`
			TotalPages    int `json:"total_pages"`
			TotalItems    int `json:"total_items"`
			ItemsPerPage  int `json:"items_per_page"`
		} `json:"pagination"`
	} `json:"data"`
}

type GroupResponse struct {
	Data models.GroupWithStats `json:"data"`
}

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type ListActivitiesResponse struct {
	Data []*models.StudyActivity `json:"data"`
}

type ActivityResponse struct {
	Data models.StudyActivity `json:"data"`
}

type CreateActivityRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type UpdateActivityRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type ListSessionsResponse struct {
	Data struct {
		Sessions   []*models.StudySession `json:"sessions"`
		Pagination struct {
			CurrentPage   int `json:"current_page"`
			TotalPages    int `json:"total_pages"`
			TotalItems    int `json:"total_items"`
			ItemsPerPage  int `json:"items_per_page"`
		} `json:"pagination"`
	} `json:"data"`
}

type SessionResponse struct {
	Data models.StudySession `json:"data"`
}

type CreateSessionRequest struct {
	GroupID         int64 `json:"group_id" binding:"required"`
	StudyActivityID int64 `json:"study_activity_id" binding:"required"`
}

type AddReviewRequest struct {
	WordID  int64 `json:"word_id" binding:"required"`
	Correct bool  `json:"correct"`
}

type ReviewResponse struct {
	Data models.WordReviewItem `json:"data"`
} 