package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"backend-go/internal/domain/models"
	"backend-go/internal/service"
	"backend-go/internal/api/handlers"
)

func setupWordTest() (*gin.Engine, *service.WordService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockRepo := service.NewMockWordRepository()
	wordService := service.NewWordService(mockRepo)
	handler := handlers.NewWordHandler(wordService)

	// Setup routes
	r.GET("/api/words", handler.ListWords)
	r.GET("/api/words/:id", handler.GetWord)
	r.POST("/api/words", handler.CreateWord)
	r.PUT("/api/words/:id", handler.UpdateWord)
	r.DELETE("/api/words/:id", handler.DeleteWord)

	return r, wordService
}

func TestWordHandler_CreateWord(t *testing.T) {
	r, _ := setupWordTest()

	tests := []struct {
		name       string
		word       models.Word
		wantStatus int
	}{
		{
			name: "valid word",
			word: models.Word{
				Kanji:   "食べる",
				Romaji:  "taberu",
				English: "to eat",
				Parts: map[string]any{
					"verb_type": "ru-verb",
					"topic":     "food",
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid word - missing kanji",
			word: models.Word{
				Romaji:  "taberu",
				English: "to eat",
				Parts:   map[string]any{},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.word)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/words", bytes.NewBuffer(body))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response map[string]models.Word
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response["data"].ID)
				assert.Equal(t, tt.word.Kanji, response["data"].Kanji)
			}
		})
	}
}

func TestWordHandler_GetWord(t *testing.T) {
	r, wordService := setupWordTest()

	// Create a test word
	word := &models.Word{
		Kanji:   "食べる",
		Romaji:  "taberu",
		English: "to eat",
		Parts: map[string]any{
			"verb_type": "ru-verb",
			"topic":     "food",
		},
	}
	err := wordService.CreateWord(nil, word)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		wordID     string
		wantStatus int
	}{
		{
			name:       "existing word",
			wordID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existent word",
			wordID:     "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid word ID",
			wordID:     "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/words/"+tt.wordID, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response map[string]models.WordWithStats
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, word.Kanji, response["data"].Kanji)
			}
		})
	}
} 