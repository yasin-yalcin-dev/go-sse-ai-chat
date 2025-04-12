/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package dto

import "github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"

// Message request and response DTOs

// CreateMessageRequest represents the request to create a new message
type CreateMessageRequest struct {
	Content string             `json:"content" binding:"required"`
	Role    models.MessageRole `json:"role"`
	Type    models.MessageType `json:"type"`
}

// MessageResponse represents the response for a message
type MessageResponse struct {
	ID        string                 `json:"id"`
	ChatID    string                 `json:"chat_id"`
	Content   string                 `json:"content"`
	Role      models.MessageRole     `json:"role"`
	Type      models.MessageType     `json:"type"`
	CreatedAt string                 `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MessageListResponse represents the response for a list of messages
type MessageListResponse struct {
	Messages   []MessageResponse `json:"messages"`
	Pagination PaginationInfo    `json:"pagination"`
}
