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
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models/dto"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/errors"
)

// CreateChat handles POST /api/v1/chats
func (h *Handler) CreateChat(c *gin.Context) {
	var req dto.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, errors.NewBadRequestError("Invalid request body", err))
		return
	}

	chat, err := h.chatService.CreateChat(c.Request.Context(), req.Title)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain model to DTO
	response := dto.ChatResponse{
		ID:           chat.ID.Hex(),
		Title:        chat.Title,
		CreatedAt:    chat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    chat.UpdatedAt.Format(time.RFC3339),
		MessageCount: chat.MessageCount,
	}

	if !chat.LastMessageAt.IsZero() {
		response.LastMessageAt = chat.LastMessageAt.Format(time.RFC3339)
	}

	respondWithJSON(c, http.StatusCreated, response)
}

// GetChat handles GET /api/v1/chats/:id
func (h *Handler) GetChat(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, errors.NewBadRequestError("Chat ID is required", nil))
		return
	}

	chat, err := h.chatService.GetChatByID(c.Request.Context(), id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain model to DTO
	response := dto.ChatResponse{
		ID:           chat.ID.Hex(),
		Title:        chat.Title,
		CreatedAt:    chat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    chat.UpdatedAt.Format(time.RFC3339),
		MessageCount: chat.MessageCount,
	}

	if !chat.LastMessageAt.IsZero() {
		response.LastMessageAt = chat.LastMessageAt.Format(time.RFC3339)
	}

	respondWithJSON(c, http.StatusOK, response)
}

// ListChats handles GET /api/v1/chats
func (h *Handler) ListChats(c *gin.Context) {
	page, pageSize := handlePagination(c, 10, 100)

	chats, total, err := h.chatService.ListChats(c.Request.Context(), page, pageSize)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain models to DTOs
	chatResponses := make([]dto.ChatResponse, len(chats))
	for i, chat := range chats {
		chatResponses[i] = dto.ChatResponse{
			ID:           chat.ID.Hex(),
			Title:        chat.Title,
			CreatedAt:    chat.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    chat.UpdatedAt.Format(time.RFC3339),
			MessageCount: chat.MessageCount,
		}

		if !chat.LastMessageAt.IsZero() {
			chatResponses[i].LastMessageAt = chat.LastMessageAt.Format(time.RFC3339)
		}
	}

	response := dto.ChatListResponse{
		Chats: chatResponses,
		Pagination: dto.PaginationInfo{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Pages:    calculateTotalPages(total, pageSize),
		},
	}

	respondWithJSON(c, http.StatusOK, response)
}

// UpdateChat handles PUT /api/v1/chats/:id
func (h *Handler) UpdateChat(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, errors.NewBadRequestError("Chat ID is required", nil))
		return
	}

	var req dto.UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, errors.NewBadRequestError("Invalid request body", err))
		return
	}

	chat, err := h.chatService.UpdateChat(c.Request.Context(), id, req.Title)
	if err != nil {
		respondWithError(c, err)
		return
	}

	// Convert domain model to DTO
	response := dto.ChatResponse{
		ID:           chat.ID.Hex(),
		Title:        chat.Title,
		CreatedAt:    chat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    chat.UpdatedAt.Format(time.RFC3339),
		MessageCount: chat.MessageCount,
	}

	if !chat.LastMessageAt.IsZero() {
		response.LastMessageAt = chat.LastMessageAt.Format(time.RFC3339)
	}

	respondWithJSON(c, http.StatusOK, response)
}

// DeleteChat handles DELETE /api/v1/chats/:id
func (h *Handler) DeleteChat(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, errors.NewBadRequestError("Chat ID is required", nil))
		return
	}

	if err := h.chatService.DeleteChat(c.Request.Context(), id); err != nil {
		respondWithError(c, err)
		return
	}

	respondWithJSON(c, http.StatusOK, dto.SuccessResponse{
		Message: "Chat deleted successfully",
	})
}
