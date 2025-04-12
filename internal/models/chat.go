/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Chat represents a chat session
type Chat struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	LastMessageAt time.Time          `bson:"last_message_at,omitempty" json:"last_message_at,omitempty"`
	MessageCount  int                `bson:"message_count" json:"message_count"`
	Active        bool               `bson:"active" json:"active"`
}

// NewChat creates a new chat with default values
func NewChat(title string) *Chat {
	now := time.Now()
	return &Chat{
		ID:           primitive.NewObjectID(),
		Title:        title,
		CreatedAt:    now,
		UpdatedAt:    now,
		MessageCount: 0,
		Active:       true,
	}
}

// BeforeSave updates the UpdatedAt field
func (c *Chat) BeforeSave() {
	c.UpdatedAt = time.Now()
}

// AddMessage updates the chat when a new message is added
func (c *Chat) AddMessage() {
	now := time.Now()
	c.LastMessageAt = now
	c.UpdatedAt = now
	c.MessageCount++
}
