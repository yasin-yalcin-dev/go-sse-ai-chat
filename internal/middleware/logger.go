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

	"github.com/gin-gonic/gin"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// LoggerMiddleware logs request information using the application's logger
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log data after request is processed
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// Create log fields
		logFields := []interface{}{
			"status", statusCode,
			"latency", latency,
			"client_ip", clientIP,
			"method", method,
			"path", path,
		}

		if errorMessage != "" {
			logFields = append(logFields, "error", errorMessage)
		}

		// Log with appropriate level based on status code
		switch {
		case statusCode >= 500:
			logger.With(logFields...).Error("Server error")
		case statusCode >= 400:
			logger.With(logFields...).Warn("Client error")
		default:
			logger.With(logFields...).Info("Request completed")
		}
	}
}
