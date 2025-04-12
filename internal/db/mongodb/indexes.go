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

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates all required indexes for the collections
func (c *DBConnection) EnsureIndexes(ctx context.Context) error {
	// Create indexes for chats collection
	if err := c.createChatIndexes(ctx); err != nil {
		return err
	}

	// Create indexes for messages collection
	if err := c.createMessageIndexes(ctx); err != nil {
		return err
	}

	logger.Info("All database indexes created successfully")
	return nil
}

// createChatIndexes creates indexes for the chats collection
func (c *DBConnection) createChatIndexes(ctx context.Context) error {
	// Define the indexes for chats collection
	chatIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "updated_at", Value: -1},
				{Key: "active", Value: 1},
			},
			Options: options.Index().SetName("updated_at_active"),
		},
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
			},
			Options: options.Index().SetName("title_text"),
		},
	}

	// Create the indexes
	_, err := c.Chats().Indexes().CreateMany(ctx, chatIndexes)
	if err != nil {
		logger.Errorf("Failed to create chat indexes: %v", err)
		return err
	}

	logger.Info("Chat indexes created successfully")
	return nil
}

// createMessageIndexes creates indexes for the messages collection
func (c *DBConnection) createMessageIndexes(ctx context.Context) error {
	// Define the indexes for messages collection
	messageIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("chat_id_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "role", Value: 1},
				{Key: "chat_id", Value: 1},
			},
			Options: options.Index().SetName("role_chat_id"),
		},
		{
			Keys: bson.D{
				{Key: "content", Value: "text"},
			},
			Options: options.Index().SetName("content_text"),
		},
	}

	// Create the indexes
	_, err := c.Messages().Indexes().CreateMany(ctx, messageIndexes)
	if err != nil {
		logger.Errorf("Failed to create message indexes: %v", err)
		return err
	}

	logger.Info("Message indexes created successfully")
	return nil
}
