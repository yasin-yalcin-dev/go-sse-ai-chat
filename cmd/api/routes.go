/*
** ** ** ** ** **

	\ \ / / \ \ / /
	 \ V /   \ V /
	  | |     | |
	  |_|     |_|
	 Yasin   Yalcin
*/
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/config"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/handlers"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/middleware"
)

// setupRouter configures the Gin router with routes and middleware
func setupRouter(cfg *config.Config) *gin.Engine {
	// Create default gin router with Logger and Recovery middleware
	router := gin.Default()

	// Add custom middleware
	router.Use(middleware.CORSMiddleware(cfg.Server.AllowedOrigins))
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Initialize handlers
	systemHandler := handlers.NewSystemHandler(cfg)

	// API v1 routes - prepare the group but don't add chat routes yet
	//apiV1 := router.Group("/api/v1")
	{
		// Chat routes will be added later
	}

	// System routes (outside of versioned API)
	systemRoutes := router.Group("/system")
	{
		systemRoutes.GET("/health", systemHandler.HealthCheck)
		systemRoutes.GET("/version", systemHandler.Version)
	}

	// Serve static files
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	return router
}
