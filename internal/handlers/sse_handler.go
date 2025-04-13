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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/services"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/sse"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/errors"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// SSEHandler handles SSE connections
type SSEHandler struct {
	broker      *sse.Broker
	chatService services.ChatService
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(broker *sse.Broker, chatService services.ChatService) *SSEHandler {
	return &SSEHandler{
		broker:      broker,
		chatService: chatService,
	}
}

// HandleStream handles streaming chat events via SSE
func (h *SSEHandler) HandleStream(c *gin.Context) {
	// Get chat ID from URL
	chatID := c.Param("id")
	if chatID == "" {
		err := errors.NewBadRequestError("Missing chat ID", nil)
		c.JSON(err.GetStatusCode(), err.ToResponse())
		return
	}

	// Verify the chat exists
	chat, err := h.chatService.GetChatByID(c.Request.Context(), chatID)
	if err != nil {
		logger.Errorf("Error fetching chat %s: %v", chatID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if chat == nil {
		err := errors.NewNotFoundError(fmt.Sprintf("Chat %s not found", chatID), nil)
		c.JSON(err.GetStatusCode(), err.ToResponse())
		return
	}

	// Create a unique client ID
	clientID := fmt.Sprintf("%s_%s", chatID, uuid.New().String())

	// Log connection attempt
	logger.Infof("SSE connection requested for chat %s, client %s", chatID, clientID)

	// Create new client
	client := sse.NewClient(clientID, c.Writer, h.broker)

	// Start listening for messages
	client.Listen()

	// The connection will be kept open until the client disconnects
	// or the context is canceled
	<-c.Request.Context().Done()
	logger.Debugf("Connection context done for client %s", clientID)
}

// GetStats returns stats about SSE connections
func (h *SSEHandler) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"clients": h.broker.GetClientCount(),
	})
}
