package handlers

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
)

func setupActivityTest() (*gin.Engine, *service.StudyActivityService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockRepo := service.NewMockStudyActivityRepository()
	sessionRepo := service.NewMockStudySessionRepository()
	activityService := service.NewStudyActivityService(mockRepo, sessionRepo)
	handler := NewStudyActivityHandler(activityService)

	// Setup routes
	r.GET("/api/study-activities", handler.ListActivities)
	r.GET("/api/study-activities/:id", handler.GetActivity)
	r.POST("/api/study-activities", handler.CreateActivity)
	r.PUT("/api/study-activities/:id", handler.UpdateActivity)
	r.DELETE("/api/study-activities/:id", handler.DeleteActivity)

	return r, activityService
}

func TestStudyActivityHandler_CreateActivity(t *testing.T) {
	r, _ := setupActivityTest()

	tests := []struct {
		name       string
		activity   models.StudyActivity
		wantStatus int
	}{
		{
			name: "valid activity",
			activity: models.StudyActivity{
				Name: "Flashcards",
				URL:  "https://example.com/flashcards",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid activity - missing name",
			activity: models.StudyActivity{
				URL: "https://example.com/flashcards",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid activity - missing URL",
			activity: models.StudyActivity{
				Name: "Flashcards",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid activity - invalid URL",
			activity: models.StudyActivity{
				Name: "Flashcards",
				URL:  "not-a-url",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.activity)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/study-activities", bytes.NewBuffer(body))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response map[string]models.StudyActivity
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response["data"].ID)
				assert.Equal(t, tt.activity.Name, response["data"].Name)
				assert.Equal(t, tt.activity.URL, response["data"].URL)
			}
		})
	}
}

func TestStudyActivityHandler_GetActivity(t *testing.T) {
	r, activityService := setupActivityTest()

	// Create a test activity
	activity := &models.StudyActivity{
		Name: "Flashcards",
		URL:  "https://example.com/flashcards",
	}
	err := activityService.CreateActivity(nil, activity)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		activityID string
		wantStatus int
	}{
		{
			name:       "existing activity",
			activityID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existent activity",
			activityID: "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid activity ID",
			activityID: "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/study-activities/"+tt.activityID, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response map[string]models.StudyActivity
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, activity.Name, response["data"].Name)
				assert.Equal(t, activity.URL, response["data"].URL)
			}
		})
	}
}

func TestStudyActivityHandler_ListActivities(t *testing.T) {
	r, activityService := setupActivityTest()

	// Create test activities
	activities := []*models.StudyActivity{
		{Name: "Flashcards", URL: "https://example.com/flashcards"},
		{Name: "Quiz", URL: "https://example.com/quiz"},
	}

	for _, activity := range activities {
		err := activityService.CreateActivity(nil, activity)
		assert.NoError(t, err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/study-activities", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]*models.StudyActivity
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, len(activities), len(response["data"]))
} 