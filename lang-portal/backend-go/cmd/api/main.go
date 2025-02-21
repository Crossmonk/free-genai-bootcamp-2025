package main

import (
	"log"
	"path/filepath"

	"backend-go/internal/api/router"
	"backend-go/internal/repository/sqlite"
	"backend-go/internal/repository/sqlite/implementations"
	"backend-go/internal/service"
	"backend-go/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	
	// Load configuration
	cfg := config.New()

	// Initialize database
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(filepath.Join(".", "migrations")); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set Gin to release mode in production
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repositories
	wordRepo := implementations.NewWordRepository(db)
	groupRepo := implementations.NewGroupRepository(db)
	activityRepo := implementations.NewStudyActivityRepository(db)
	sessionRepo := implementations.NewStudySessionRepository(db)

	// Initialize services
	wordService := service.NewWordService(wordRepo)
	groupService := service.NewGroupService(groupRepo)
	activityService := service.NewStudyActivityService(activityRepo, sessionRepo)
	sessionService := service.NewStudySessionService(sessionRepo, groupRepo)

	// Initialize router with services
	r := router.SetupRouter(wordService, groupService, activityService, sessionService)

	// Basic health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start server
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 