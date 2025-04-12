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

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/db/mongodb"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MessageRepository implements the MessageRepository interface
type MessageRepository struct {
	db *mongodb.DBConnection
}

// NewMessageRepository creates a new MongoDB message repository
func NewMessageRepository(db *mongodb.DBConnection) repository.MessageRepository {
	return &MessageRepository{db: db}
}

// Create inserts a new message into the database
func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	_, err := r.db.Messages().InsertOne(ctx, message)
	return err
}

// FindByID retrieves a message by its ID
func (r *MessageRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Message, error) {
	var message models.Message
	err := r.db.Messages().FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Message not found
		}
		return nil, err
	}
	return &message, nil
}

// FindByChatID retrieves messages for a specific chat with pagination
func (r *MessageRepository) FindByChatID(ctx context.Context, chatID primitive.ObjectID, limit, offset int) ([]*models.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}}). // Sort by created_at descending (newest first)
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.db.Messages().Find(ctx, bson.M{"chat_id": chatID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverse the order to have chronological order (oldest first)
	// This is often better for displaying chat messages
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// CountByChatID counts the number of messages in a chat
func (r *MessageRepository) CountByChatID(ctx context.Context, chatID primitive.ObjectID) (int64, error) {
	return r.db.Messages().CountDocuments(ctx, bson.M{"chat_id": chatID})
}

// Delete removes a message
func (r *MessageRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.db.Messages().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// DeleteByChatID removes all messages for a specific chat
func (r *MessageRepository) DeleteByChatID(ctx context.Context, chatID primitive.ObjectID) error {
	_, err := r.db.Messages().DeleteMany(ctx, bson.M{"chat_id": chatID})
	return err
}
