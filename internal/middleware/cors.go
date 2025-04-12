/*
** ** ** ** ** **

	\ \ / / \ \ / /
	 \ V /   \ V /
	  | |     | |
	  |_|     |_|
	 Yasin   Yalcin
*/

package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a middleware for handling Cross-Origin Resource Sharing (CORS)
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	// Default CORS configuration
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Handle wildcard case
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = allowedOrigins
	}

	return cors.New(corsConfig)
}
