package models

import "time"

type StudySessionWithStats struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Embedded study activity info
	StudyActivity struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"study_activity"`

	// Session statistics
	Stats struct {
		TotalReviews    int     `json:"total_reviews"`
		CorrectReviews  int     `json:"correct_reviews"`
		Accuracy        float64 `json:"accuracy"`
		DurationMinutes int     `json:"duration_minutes"`
	} `json:"stats"`
}

type StudyProgress struct {
	TotalSessions  int `json:"total_sessions"`
	TotalReviews   int `json:"total_reviews"`
	CorrectReviews int `json:"correct_reviews"`
	TimeRange struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	} `json:"time_range"`
	DailyStats []DailyStats `json:"daily_stats"`
}

type DailyStats struct {
	Date           string  `json:"date"`
	TotalSessions  int     `json:"total_sessions"`
	TotalReviews   int     `json:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews"`
	Accuracy       float64 `json:"accuracy"`
}

type QuickStats struct {
	TotalSessions int     `json:"total_sessions"`
	TotalReviews  int     `json:"total_reviews"`
	Accuracy      float64 `json:"accuracy"`
}

type WeekStats struct {
	TotalSessions  int     `json:"total_sessions"`
	TotalReviews   int     `json:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews"`
	Accuracy       float64 `json:"accuracy"`
}

type AllTimeStats struct {
	TotalSessions  int     `json:"total_sessions"`
	TotalReviews   int     `json:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews"`
	Accuracy       float64 `json:"accuracy"`
}

type ListSessionsResult struct {
	Sessions    []*StudySession
	TotalItems  int
	CurrentPage int
	TotalPages  int
} 