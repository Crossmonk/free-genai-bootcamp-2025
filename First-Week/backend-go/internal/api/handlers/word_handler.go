package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend-go/internal/domain/models"
	"backend-go/internal/service"
)

type WordHandler struct {
	wordService *service.WordService
}

func NewWordHandler(wordService *service.WordService) *WordHandler {
	return &WordHandler{
		wordService: wordService,
	}
}

// ListWords godoc
// @Summary List words with pagination and sorting
// @Description Get a paginated list of words with their review statistics
// @Tags words
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Param sort_by query string false "Sort field (kanji, romaji, english, correct_count, wrong_count)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} ListWordsResponse
// @Router /api/words [get]
func (h *WordHandler) ListWords(c *gin.Context) {
	params := service.ListWordsParams{
		Page:     parseInt(c.Query("page"), 1),
		PageSize: parseInt(c.Query("page_size"), 10),
		SortBy:   c.Query("sort_by"),
		Order:    c.Query("order"),
	}

	result, err := h.wordService.ListWords(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"words": result.Words,
			"pagination": gin.H{
				"current_page": result.CurrentPage,
				"total_pages": result.TotalPages,
				"total_items": result.TotalItems,
				"items_per_page": params.PageSize,
			},
		},
	})
}

// GetWord godoc
// @Summary Get a word by ID
// @Description Get a single word by its ID with review statistics
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Success 200 {object} WordResponse
// @Router /api/words/{id} [get]
func (h *WordHandler) GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word ID"})
		return
	}

	word, err := h.wordService.GetWordWithStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}

// CreateWord godoc
// @Summary Create a new word
// @Description Create a new vocabulary word
// @Tags words
// @Accept json
// @Produce json
// @Param word body CreateWordRequest true "Word object"
// @Success 201 {object} WordResponse
// @Router /api/words [post]
func (h *WordHandler) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.wordService.CreateWord(c.Request.Context(), &word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": word,
	})
}

// UpdateWord godoc
// @Summary Update a word
// @Description Update an existing vocabulary word
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Param word body UpdateWordRequest true "Word object"
// @Success 200 {object} WordResponse
// @Router /api/words/{id} [put]
func (h *WordHandler) UpdateWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word ID"})
		return
	}

	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	word.ID = id

	if err := h.wordService.UpdateWord(c.Request.Context(), &word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}

// DeleteWord godoc
// @Summary Delete a word
// @Description Delete a vocabulary word
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Success 204 "No Content"
// @Router /api/words/{id} [delete]
func (h *WordHandler) DeleteWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word ID"})
		return
	}

	if err := h.wordService.DeleteWord(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseInt(str string, defaultValue int) int {
	if str == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}