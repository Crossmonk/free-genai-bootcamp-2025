package models

type Word struct {
	ID      int64           `json:"id"`
	Kanji   string         `json:"kanji"`
	Romaji  string         `json:"romaji"`
	English string         `json:"english"`
	Parts   map[string]any `json:"parts"`
}

type WordStats struct {
	CorrectCount int     `json:"correct_count"`
	WrongCount   int     `json:"wrong_count"`
	Accuracy     float64 `json:"accuracy"`
}

type WordWithStats struct {
	ID      int64           `json:"id"`
	Kanji   string         `json:"kanji"`
	Romaji  string         `json:"romaji"`
	English string         `json:"english"`
	Parts   map[string]any `json:"parts"`
	Stats   struct {
		CorrectCount int     `json:"correct_count"`
		WrongCount   int     `json:"wrong_count"`
		Accuracy     float64 `json:"accuracy"`
	} `json:"stats"`
} 