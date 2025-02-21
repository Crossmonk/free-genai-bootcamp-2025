package models

import "time"

type Group struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	WordsCount    int        `json:"words_count"`
	LastStudiedAt *time.Time `json:"last_studied_at,omitempty"`
}

type GroupStats struct {
	TotalReviews   int     `json:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews"`
	Accuracy       float64 `json:"accuracy"`
}

type GroupWithStats struct {
	Group
	Stats GroupStats `json:"stats"`
} 