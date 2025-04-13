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
	"strings"

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

	// Get client ID from query param (for reconnection) or generate new one
	clientID := c.Query("client_id")

	// If no client ID or invalid format, generate a new one
	if clientID == "" || !strings.HasPrefix(clientID, chatID+"_") {
		clientID = fmt.Sprintf("%s_%s", chatID, uuid.New().String())
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

	// Get last event ID for replay (if client is reconnecting)
	lastEventID := c.GetHeader("Last-Event-ID")
	isReconnection := lastEventID != ""

	// Log connection attempt
	if isReconnection {
		logger.Infof("SSE reconnection requested for chat %s, client %s, last event %s",
			chatID, clientID, lastEventID)
	} else {
		logger.Infof("New SSE connection requested for chat %s, client %s", chatID, clientID)
	}

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
	// Get stats per chat if chat ID is provided
	chatID := c.Query("chat_id")
	if chatID != "" {
		clientCount := h.broker.GetClientsInChat(chatID)
		c.JSON(http.StatusOK, gin.H{
			"chat_id": chatID,
			"clients": clientCount,
		})
		return
	}

	// Otherwise return total stats
	c.JSON(http.StatusOK, gin.H{
		"total_clients": h.broker.GetClientCount(),
	})
}
