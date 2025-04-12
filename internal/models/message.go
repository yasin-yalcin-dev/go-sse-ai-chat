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

// MessageRole represents the role of a message sender
type MessageRole string

// Message roles
const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// MessageType represents the type of message content
type MessageType string

// Message types
const (
	TypeText  MessageType = "text"
	TypeImage MessageType = "image"
	TypeCode  MessageType = "code"
)

// Message represents a message in a chat
type Message struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ChatID    primitive.ObjectID     `bson:"chat_id" json:"chat_id"`
	Content   string                 `bson:"content" json:"content"`
	Role      MessageRole            `bson:"role" json:"role"`
	Type      MessageType            `bson:"type" json:"type"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// NewMessage creates a new message with default values
func NewMessage(chatID primitive.ObjectID, content string, role MessageRole, msgType MessageType) *Message {
	return &Message{
		ID:        primitive.NewObjectID(),
		ChatID:    chatID,
		Content:   content,
		Role:      role,
		Type:      msgType,
		CreatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// SetMetadata adds or updates a metadata entry
func (m *Message) SetMetadata(key string, value interface{}) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	m.Metadata[key] = value
}

// GetMetadata retrieves a metadata value
func (m *Message) GetMetadata(key string) (interface{}, bool) {
	if m.Metadata == nil {
		return nil, false
	}
	val, ok := m.Metadata[key]
	return val, ok
}
