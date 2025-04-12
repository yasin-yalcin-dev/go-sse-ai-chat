/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package repository

import (
	"context"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRepository defines the interface for chat data access
type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Chat, error)
	FindAll(ctx context.Context, limit, offset int) ([]*models.Chat, error)
	Update(ctx context.Context, chat *models.Chat) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	IncrementMessageCount(ctx context.Context, id primitive.ObjectID) error
	CountAll(ctx context.Context) (int64, error)
}

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Message, error)
	FindByChatID(ctx context.Context, chatID primitive.ObjectID, limit, offset int) ([]*models.Message, error)
	CountByChatID(ctx context.Context, chatID primitive.ObjectID) (int64, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteByChatID(ctx context.Context, chatID primitive.ObjectID) error
}
