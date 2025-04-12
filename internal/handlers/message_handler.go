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
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models/dto"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/errors"
)

// GetMessages handles GET /api/v1/chats/:id/messages
func (h *Handler) GetMessages(c *gin.Context) {
	chatID := c.Param("id")
	if chatID == "" {
		respondWithError(c, errors.NewBadRequestError("Chat ID is required", nil))
		return
	}

	page, pageSize := handlePagination(c, 20, 100)

	messages, total, err := h.messageService.GetChatMessages(c.Request.Context(), chatID, page, pageSize)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain models to DTOs
	messageResponses := make([]dto.MessageResponse, len(messages))
	for i, message := range messages {
		messageResponses[i] = dto.MessageResponse{
			ID:        message.ID.Hex(),
			ChatID:    message.ChatID.Hex(),
			Content:   message.Content,
			Role:      message.Role,
			Type:      message.Type,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
			Metadata:  message.Metadata,
		}
	}

	response := dto.MessageListResponse{
		Messages: messageResponses,
		Pagination: dto.PaginationInfo{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Pages:    calculateTotalPages(total, pageSize),
		},
	}

	respondWithJSON(c, http.StatusOK, response)
}

// CreateMessage handles POST /api/v1/chats/:id/messages
func (h *Handler) CreateMessage(c *gin.Context) {
	chatID := c.Param("id")
	if chatID == "" {
		respondWithError(c, errors.NewBadRequestError("Chat ID is required", nil))
		return
	}

	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, errors.NewBadRequestError("Invalid request body", err))
		return
	}

	// Set default values if not provided
	if req.Role == "" {
		req.Role = models.RoleUser
	}

	if req.Type == "" {
		req.Type = models.TypeText
	}

	// Validate role
	if req.Role != models.RoleUser && req.Role != models.RoleAssistant && req.Role != models.RoleSystem {
		respondWithError(c, errors.NewValidationError("Invalid role", nil))
		return
	}

	// Validate type
	if req.Type != models.TypeText && req.Type != models.TypeImage && req.Type != models.TypeCode {
		respondWithError(c, errors.NewValidationError("Invalid message type", nil))
		return
	}

	message, err := h.messageService.CreateMessage(
		c.Request.Context(),
		chatID,
		req.Content,
		req.Role,
		req.Type,
	)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain model to DTO
	response := dto.MessageResponse{
		ID:        message.ID.Hex(),
		ChatID:    message.ChatID.Hex(),
		Content:   message.Content,
		Role:      message.Role,
		Type:      message.Type,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
		Metadata:  message.Metadata,
	}

	respondWithJSON(c, http.StatusCreated, response)
}

// GetMessage handles GET /api/v1/messages/:id
func (h *Handler) GetMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, errors.NewBadRequestError("Message ID is required", nil))
		return
	}

	message, err := h.messageService.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain model to DTO
	response := dto.MessageResponse{
		ID:        message.ID.Hex(),
		ChatID:    message.ChatID.Hex(),
		Content:   message.Content,
		Role:      message.Role,
		Type:      message.Type,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
		Metadata:  message.Metadata,
	}

	respondWithJSON(c, http.StatusOK, response)
}

// DeleteMessage handles DELETE /api/v1/messages/:id
func (h *Handler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, errors.NewBadRequestError("Message ID is required", nil))
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), id); err != nil {
		respondWithError(c, err)
		return
	}

	respondWithJSON(c, http.StatusOK, dto.SuccessResponse{
		Message: "Message deleted successfully",
	})
}
