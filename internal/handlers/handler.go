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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models/dto"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/services"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/errors"
)

// Handler contains services for all handlers
type Handler struct {
	chatService    services.ChatService
	messageService services.MessageService
}

// NewHandler creates a new handler with all required services
func NewHandler(chatService services.ChatService, messageService services.MessageService) *Handler {
	return &Handler{
		chatService:    chatService,
		messageService: messageService,
	}
}

// respondWithError sends a JSON error response
func respondWithError(c *gin.Context, err error) {
	// Check if it's our application error type
	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		response := dto.ErrorResponse{
			Error: dto.ErrorDetails{
				Code:    appErr.Code,
				Message: appErr.Error(),
			},
		}

		// Add details if available
		if len(appErr.Context) > 0 {
			response.Error.Details = "See error context for more details"
		}

		c.JSON(appErr.GetStatusCode(), response)
		return
	}

	// If it's a generic error, return a 500 Internal Server Error
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Error: dto.ErrorDetails{
			Message: err.Error(),
		},
	})
}

// respondWithJSON sends a JSON response
func respondWithJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// handlePagination extracts pagination parameters from the request
func handlePagination(c *gin.Context, defaultPageSize, maxPageSize int) (page, pageSize int) {
	page = 1
	pageSize = defaultPageSize

	if pageStr := c.Query("page"); pageStr != "" {
		if pageInt, err := strconv.Atoi(pageStr); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSizeInt, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeInt > 0 {
			pageSize = pageSizeInt
		}
	}

	// Ensure page size doesn't exceed maximum
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}

// calculateTotalPages calculates the total number of pages
func calculateTotalPages(total int64, pageSize int) int64 {
	return (total + int64(pageSize) - 1) / int64(pageSize)
}
