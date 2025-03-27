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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// setupRouter configures the Gin router with basic routes
func setupRouter() *gin.Engine {
	// Create default gin router with Logger and Recovery middleware
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Version endpoint
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":    "0.1.0",
			"build_time": time.Now().Format(time.RFC3339),
		})
	})

	// Serve static files
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	return router
}
