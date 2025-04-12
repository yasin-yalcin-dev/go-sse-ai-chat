/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/config"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// SystemHandler handles system-related endpoints
type SystemHandler struct {
	config *config.Config
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(cfg *config.Config) *SystemHandler {
	return &SystemHandler{
		config: cfg,
	}
}

// HealthCheck handles the health check endpoint
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	// Create health status response
	response := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"services":  make(map[string]string),
	}

	// Add version information
	response["version"] = "0.1.0"

	// TODO: Add status of dependencies (MongoDB, etc.)
	// For now, just mark them as "unknown"
	response["services"].(map[string]string)["database"] = "unknown"
	response["services"].(map[string]string)["ai_provider"] = "unknown"

	logger.Debug("Health check executed")
	c.JSON(http.StatusOK, response)
}

// Version returns the current API version
func (h *SystemHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    "0.1.0",
		"build_time": "2025-04-12T00:00:00Z", // This should be set during build process
		"go_version": "1.20",
	})
}
