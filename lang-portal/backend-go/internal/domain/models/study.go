package models

import "time"

type StudyActivity struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type StudySessionStats struct {
	TotalReviews    int     `json:"total_reviews"`
	CorrectReviews  int     `json:"correct_reviews"`
	Accuracy        float64 `json:"accuracy"`
	DurationMinutes int     `json:"duration_minutes"`
} 