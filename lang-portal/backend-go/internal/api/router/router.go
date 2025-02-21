package router

import (
	"backend-go/internal/api/handlers"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and configures all routes
func SetupRouter(
	wordService *service.WordService,
	groupService *service.GroupService,
	activityService *service.StudyActivityService,
	sessionService *service.StudySessionService,
) *gin.Engine {
	router := gin.Default()

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordService)
	groupHandler := handlers.NewGroupHandler(groupService)
	activityHandler := handlers.NewStudyActivityHandler(activityService)
	sessionHandler := handlers.NewStudySessionHandler(sessionService)

	// API group
	api := router.Group("/api")
	{
		// Words routes
		api.GET("/words", wordHandler.ListWords)
		api.GET("/words/:id", wordHandler.GetWord)

		// Groups routes
		api.GET("/groups", groupHandler.ListGroups)
		api.GET("/group/:id", groupHandler.GetGroup)
		api.GET("/group/:id/words", groupHandler.GetGroupWords)
		api.GET("/group/:id/study_sessions", groupHandler.GetGroupStudySessions)

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/last_study_session", sessionHandler.GetLastStudySession)
			dashboard.GET("/study_progress", sessionHandler.GetStudyProgress)
			dashboard.GET("/quick_stats", sessionHandler.GetQuickStats)
		}

		// Study activities routes
		api.GET("/study_activity/:id", activityHandler.GetActivity)
		api.GET("/study_activity/:id/study_sessions", activityHandler.ListSessions)
		api.POST("/study_activities", activityHandler.CreateActivity)

		// Study sessions routes
		api.GET("/study_session/:id/words", sessionHandler.GetSessionWords)
		api.POST("/study_sessions", sessionHandler.CreateSession)
		api.POST("/study_sessions/:id/review", sessionHandler.AddReview)

		// Settings routes
		settings := api.Group("/settings")
		{
			settings.POST("/full_reset", sessionHandler.FullReset)
			settings.POST("/load_seed_data", sessionHandler.LoadSeedData)
		}
	}

	return router
}