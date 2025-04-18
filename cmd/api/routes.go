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
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/config"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/db/mongodb"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/handlers"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/middleware"
	repo "github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/repository/mongodb"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/services"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/sse"
)

// setupRouter configures the Gin router with routes and middleware
func setupRouter(cfg *config.Config) *gin.Engine {
	// Create default gin router with Logger and Recovery middleware
	router := gin.Default()

	// Add custom middleware
	router.Use(middleware.CORSMiddleware(cfg.Server.AllowedOrigins))
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	// Initialize database connection
	db, err := mongodb.New(&cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Initialize handlers
	systemHandler := handlers.NewSystemHandler(cfg)

	// Initialize repositories
	chatRepo := repo.NewChatRepository(db)
	messageRepo := repo.NewMessageRepository(db)

	// Initialize services
	chatService := services.NewChatService(chatRepo, messageRepo)
	messageService := services.NewMessageService(messageRepo, chatRepo)

	// Initialize SSE broker
	broker := sse.NewBroker(cfg.SSE.MaxClients, cfg.SSE.KeepaliveInterval)

	// Start broker in a goroutine
	go broker.Start(context.Background())

	// Initialize handlers
	handler := handlers.NewHandler(chatService, messageService)
	sseHandler := handlers.NewSSEHandler(broker, chatService)

	apiV1 := router.Group("/api/v1")
	{
		// Chat routes
		chats := apiV1.Group("/chats")
		{
			chats.GET("", handler.ListChats)
			chats.POST("", handler.CreateChat)
			chats.GET("/:id", handler.GetChat)
			chats.PUT("/:id", handler.UpdateChat)
			chats.DELETE("/:id", handler.DeleteChat)

			// Message routes (nested under chat)
			chats.GET("/:id/messages", handler.GetMessages)
			chats.POST("/:id/messages", handler.CreateMessage)

			// SSE streaming route
			chats.GET("/:id/stream", sseHandler.HandleStream)
		}

		// Individual message routes
		messages := apiV1.Group("/messages")
		{
			messages.GET("/:id", handler.GetMessage)
			messages.DELETE("/:id", handler.DeleteMessage)
		}
		// SSE stats (for monitoring)
		sse := apiV1.Group("/sse")
		{
			sse.GET("/stats", sseHandler.GetStats)
		}
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
