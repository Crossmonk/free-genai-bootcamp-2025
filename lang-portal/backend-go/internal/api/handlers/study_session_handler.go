package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/responses"
	"backend-go/internal/service"
)

type StudySessionHandler struct {
	sessionService *service.StudySessionService
}

func NewStudySessionHandler(sessionService *service.StudySessionService) *StudySessionHandler {
	return &StudySessionHandler{
		sessionService: sessionService,
	}
}

// CreateSession godoc
// @Summary Create a new study session
// @Description Create a new study session for a group
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param session body CreateSessionRequest true "Session object"
// @Success 201 {object} SessionResponse
// @Router /api/study-sessions [post]
func (h *StudySessionHandler) CreateSession(c *gin.Context) {
	var params service.CreateSessionParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	session, err := h.sessionService.CreateSession(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": session,
	})
}

// GetSession godoc
// @Summary Get a study session by ID
// @Description Get a single study session by its ID
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} SessionResponse
// @Router /api/study-sessions/{id} [get]
func (h *StudySessionHandler) GetSession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": session,
	})
}

// ListSessions godoc
// @Summary List study sessions for a group
// @Description Get a paginated list of study sessions for a group
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param group_id query int true "Group ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Success 200 {object} ListSessionsResponse
// @Router /api/study-sessions [get]
func (h *StudySessionHandler) ListSessions(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	params := service.ListSessionsParams{
		GroupID:  groupID,
		Page:     parseInt(c.Query("page"), 1),
		PageSize: parseInt(c.Query("page_size"), 10),
	}

	result, err := h.sessionService.ListSessions(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"sessions": result.Sessions,
			"pagination": gin.H{
				"current_page":   result.CurrentPage,
				"total_pages":    result.TotalPages,
				"total_items":    result.TotalItems,
				"items_per_page": params.PageSize,
			},
		},
	})
}

// AddReview godoc
// @Summary Add a word review to a session
// @Description Record a word review result in a study session
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Param review body AddReviewRequest true "Review object"
// @Success 201 {object} ReviewResponse
// @Router /api/study-sessions/{id}/reviews [post]
func (h *StudySessionHandler) AddReview(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req struct {
		WordID  int64 `json:"word_id" binding:"required"`
		Correct bool  `json:"correct"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	review, err := h.sessionService.AddReview(c.Request.Context(), sessionID, req.WordID, req.Correct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": review,
	})
}

// GetLastStudySession handles GET /api/dashboard/last_study_session
func (h *StudySessionHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.sessionService.GetLastStudySession(c.Request.Context())
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch last study session")
		return
	}

	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": gin.H{
			"session": session,
		},
	})
}

// GetStudyProgress handles GET /api/dashboard/study_progress
func (h *StudySessionHandler) GetStudyProgress(c *gin.Context) {
	// Get time range from query params, default to last 7 days
	days := parseInt(c.DefaultQuery("days", "7"), 7)
	
	progress, err := h.sessionService.GetStudyProgress(c.Request.Context(), days)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch study progress")
		return
	}

	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": progress,
	})
}

// GetQuickStats handles GET /api/dashboard/quick_stats
func (h *StudySessionHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.sessionService.GetQuickStats(c.Request.Context())
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch quick stats")
		return
	}

	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": stats,
	})
}

// GetSessionWords handles GET /study_session/:id/words
func (h *StudySessionHandler) GetSessionWords(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	words, err := h.sessionService.GetSessionWords(c.Request.Context(), sessionID)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch session words")
		return
	}

	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": gin.H{
			"words": words,
		},
	})
}

// FullReset handles the request to reset all study data
func (h *StudySessionHandler) FullReset(c *gin.Context) {
	err := h.sessionService.FullReset(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset data: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Data reset successful"})
}

// LoadSeedData handles the request to load initial seed data
func (h *StudySessionHandler) LoadSeedData(c *gin.Context) {
	err := h.sessionService.LoadSeedData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load seed data: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Seed data loaded successfully"})
} 