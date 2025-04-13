/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package services

import (
	"context"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/ai"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
)

// ChatService defines operations for managing chat sessions
type ChatService interface {
	CreateChat(ctx context.Context, title string) (*models.Chat, error)
	GetChatByID(ctx context.Context, id string) (*models.Chat, error)
	ListChats(ctx context.Context, page, pageSize int) ([]*models.Chat, int64, error)
	UpdateChat(ctx context.Context, id string, title string) (*models.Chat, error)
	DeleteChat(ctx context.Context, id string) error
}

// MessageService defines operations for managing messages
type MessageService interface {
	CreateMessage(ctx context.Context, chatID string, content string, role models.MessageRole, msgType models.MessageType) (*models.Message, error)
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	GetChatMessages(ctx context.Context, chatID string, page, pageSize int) ([]*models.Message, int64, error)
	DeleteMessage(ctx context.Context, id string) error
	GetMessageContext(ctx context.Context, chatID string, maxMessages int) ([]ai.AIMessage, error)
}
