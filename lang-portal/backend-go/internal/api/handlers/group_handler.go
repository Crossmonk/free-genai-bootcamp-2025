package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/domain/models"
	"backend-go/internal/responses"
	"backend-go/internal/service"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// ListGroups godoc
// @Summary List groups with pagination and sorting
// @Description Get a paginated list of groups with their statistics
// @Tags groups
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Param sort_by query string false "Sort field (name, words_count)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} ListGroupsResponse
// @Router /api/groups [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	params := service.ListGroupsParams{
		Page:     parseInt(c.Query("page"), 1),
		PageSize: parseInt(c.Query("page_size"), 10),
		SortBy:   c.Query("sort_by"),
		Order:    c.Query("order"),
	}

	result, err := h.groupService.ListGroups(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"groups": result.Groups,
			"pagination": gin.H{
				"current_page":   result.CurrentPage,
				"total_pages":    result.TotalPages,
				"total_items":    result.TotalItems,
				"items_per_page": params.PageSize,
			},
		},
	})
}

// GetGroup godoc
// @Summary Get a group by ID
// @Description Get a single group by its ID with statistics
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} GroupResponse
// @Router /api/groups/{id} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	group, err := h.groupService.GetGroupWithStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": group,
	})
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new word group
// @Tags groups
// @Accept json
// @Produce json
// @Param group body CreateGroupRequest true "Group object"
// @Success 201 {object} GroupResponse
// @Router /api/groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.groupService.CreateGroup(c.Request.Context(), &group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": group,
	})
}

// UpdateGroup godoc
// @Summary Update a group
// @Description Update an existing word group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param group body UpdateGroupRequest true "Group object"
// @Success 200 {object} GroupResponse
// @Router /api/groups/{id} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	group.ID = id

	if err := h.groupService.UpdateGroup(c.Request.Context(), &group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": group,
	})
}

// DeleteGroup godoc
// @Summary Delete a group
// @Description Delete a word group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 204 "No Content"
// @Router /api/groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	if err := h.groupService.DeleteGroup(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetGroupWords handles GET /group/:id/words
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	// Get group ID from URL parameter
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Get pagination parameters
	page := c.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)
	sortBy := c.DefaultQuery("sort_by", "kanji")
	order := c.DefaultQuery("order", "asc")

	// Get words from service
	words, pagination, err := h.groupService.GetGroupWords(c.Request.Context(), groupID, pageNum, sortBy, order)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch group words")
		return
	}

	// Return response
	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": gin.H{
			"group": gin.H{
				"id":   groupID,
				"name": words.GroupName,
			},
			"words":      words.Words,
			"pagination": pagination,
		},
	})
}

// GetGroupStudySessions handles GET /group/:id/study_sessions
func (h *GroupHandler) GetGroupStudySessions(c *gin.Context) {
	// Get group ID from URL parameter
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Get pagination parameters
	page := c.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)
	pageSize := c.DefaultQuery("page_size", "10")
	pageSizeNum, _ := strconv.Atoi(pageSize)

	// Get study sessions from service
	sessions, err := h.groupService.GetGroupStudySessions(c.Request.Context(), groupID, pageNum, pageSizeNum)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch group study sessions")
		return
	}

	// Return response
	responses.SuccessResponse(c, http.StatusOK, gin.H{
		"data": gin.H{
			"group": gin.H{
				"id":   groupID,
				"name": sessions.GroupName,
			},
			"study_sessions": sessions.StudySessions,
			"pagination": gin.H{
				"current_page":   sessions.CurrentPage,
				"total_pages":    sessions.TotalPages,
				"total_items":    sessions.TotalItems,
				"items_per_page": pageSizeNum,
			},
		},
	})
} 