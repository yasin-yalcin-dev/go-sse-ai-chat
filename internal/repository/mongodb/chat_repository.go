/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/db/mongodb"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatRepository implements the ChatRepository interface
type ChatRepository struct {
	db *mongodb.DBConnection
}

// NewChatRepository creates a new MongoDB chat repository
func NewChatRepository(db *mongodb.DBConnection) repository.ChatRepository {
	return &ChatRepository{db: db}
}

// Create inserts a new chat into the database
func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	chat.BeforeSave()
	_, err := r.db.Chats().InsertOne(ctx, chat)
	return err
}

// FindByID retrieves a chat by its ID
func (r *ChatRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Chats().FindOne(ctx, bson.M{"_id": id}).Decode(&chat)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Chat not found
		}
		return nil, err
	}
	return &chat, nil
}

// FindAll retrieves all chats with pagination
func (r *ChatRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Chat, error) {
	opts := options.Find().
		SetSort(bson.D{{"updated_at", -1}}). // Sort by updated_at descending
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.db.Chats().Find(ctx, bson.M{"active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []*models.Chat
	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}
	return chats, nil
}

// Update updates an existing chat
func (r *ChatRepository) Update(ctx context.Context, chat *models.Chat) error {
	chat.BeforeSave()
	result, err := r.db.Chats().ReplaceOne(ctx, bson.M{"_id": chat.ID}, chat)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Delete marks a chat as inactive (soft delete)
func (r *ChatRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"active":     false,
			"updated_at": time.Now(),
		},
	}
	result, err := r.db.Chats().UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// IncrementMessageCount increases the message count of a chat
func (r *ChatRepository) IncrementMessageCount(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"message_count": 1},
		"$set": bson.M{
			"last_message_at": now,
			"updated_at":      now,
		},
	}
	result, err := r.db.Chats().UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// CountAll returns the total number of active chats
func (r *ChatRepository) CountAll(ctx context.Context) (int64, error) {
	return r.db.Chats().CountDocuments(ctx, bson.M{"active": true})
}
