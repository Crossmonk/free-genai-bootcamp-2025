package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/domain/models"
	"backend-go/internal/responses"
	"backend-go/internal/service"
)

type StudyActivityHandler struct {
	activityService *service.StudyActivityService
}

func NewStudyActivityHandler(activityService *service.StudyActivityService) *StudyActivityHandler {
	return &StudyActivityHandler{
		activityService: activityService,
	}
}

// ListActivities godoc
// @Summary List study activities
// @Description Get a list of study activities
// @Tags study-activities
// @Accept json
// @Produce json
// @Success 200 {object} ListActivitiesResponse
// @Router /api/study-activities [get]
func (h *StudyActivityHandler) ListActivities(c *gin.Context) {
	activities, err := h.activityService.ListActivities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": activities,
	})
}

// GetActivity godoc
// @Summary Get a study activity by ID
// @Description Get a single study activity by its ID
// @Tags study-activities
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} ActivityResponse
// @Router /api/study-activities/{id} [get]
func (h *StudyActivityHandler) GetActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity ID"})
		return
	}

	activity, err := h.activityService.GetActivity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": activity,
	})
}

// CreateActivity godoc
// @Summary Create a new study activity
// @Description Create a new study activity
// @Tags study-activities
// @Accept json
// @Produce json
// @Param activity body CreateActivityRequest true "Activity object"
// @Success 201 {object} ActivityResponse
// @Router /api/study-activities [post]
func (h *StudyActivityHandler) CreateActivity(c *gin.Context) {
	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.activityService.CreateActivity(c.Request.Context(), &activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": activity,
	})
}

// UpdateActivity godoc
// @Summary Update a study activity
// @Description Update an existing study activity
// @Tags study-activities
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Param activity body UpdateActivityRequest true "Activity object"
// @Success 200 {object} ActivityResponse
// @Router /api/study-activities/{id} [put]
func (h *StudyActivityHandler) UpdateActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity ID"})
		return
	}

	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	activity.ID = id

	if err := h.activityService.UpdateActivity(c.Request.Context(), &activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": activity,
	})
}

// DeleteActivity godoc
// @Summary Delete a study activity
// @Description Delete a study activity
// @Tags study-activities
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 204 "No Content"
// @Router /api/study-activities/{id} [delete]
func (h *StudyActivityHandler) DeleteActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity ID"})
		return
	}

	if err := h.activityService.DeleteActivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSessions handles GET /api/study_activity/:id/study_sessions
func (h *StudyActivityHandler) ListSessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	page := parseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := parseInt(c.DefaultQuery("page_size", "10"), 10)

	sessions, err := h.activityService.ListActivitySessions(c.Request.Context(), id, page, pageSize)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch activity sessions")
		return
	}

	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": sessions,
	})
} 