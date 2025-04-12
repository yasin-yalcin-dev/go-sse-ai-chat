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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/errors"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// ErrorHandlerMiddleware is a middleware that catches and formats errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// If no errors in context, return
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Check if it's our application error type
		var appError *errors.AppError
		statusCode := http.StatusInternalServerError
		errorResponse := gin.H{
			"error": gin.H{
				"message": "Internal Server Error",
			},
		}

		if errors.As(err, &appError) {
			// Use app error information
			statusCode = appError.GetStatusCode()
			errorResponse = appError.ToResponse()

			// Log error details
			logger.With(
				"status_code", statusCode,
				"error_code", appError.Code,
				"path", c.Request.URL.Path,
			).Error(appError.Error())
		} else {
			// Generic error, don't leak internal details in production
			logger.With(
				"status_code", statusCode,
				"path", c.Request.URL.Path,
			).Error(err.Error())
		}

		// Add request ID to response if available
		if requestID, exists := c.Get("RequestID"); exists {
			errorResponse["request_id"] = requestID
		}

		// Respond with JSON
		c.JSON(statusCode, errorResponse)
		c.Abort()
	}
}
