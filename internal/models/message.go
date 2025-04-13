/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

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

// ProcessingStatus represents the processing status of a message
type ProcessingStatus string

// Processing statuses
const (
	StatusPending  ProcessingStatus = "pending"
	StatusSending  ProcessingStatus = "sending"
	StatusComplete ProcessingStatus = "complete"
	StatusError    ProcessingStatus = "error"
)

// Known metadata keys
const (
	MetaClientIP        = "client_ip"
	MetaUserAgent       = "user_agent"
	MetaProcessingTime  = "processing_time_ms"
	MetaErrorDetail     = "error_detail"
	MetaTokenCount      = "token_count"
	MetaModelName       = "model_name"
	MetaCodeLanguage    = "code_language"
	MetaIsEdited        = "is_edited"
	MetaOriginalContent = "original_content"
	MetaReplyTo         = "reply_to"
	MetaIsDraft         = "is_draft"
)

// Message represents a message in a chat
type Message struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ChatID           primitive.ObjectID     `bson:"chat_id" json:"chat_id"`
	Content          string                 `bson:"content" json:"content"`
	Role             MessageRole            `bson:"role" json:"role"`
	Type             MessageType            `bson:"type" json:"type"`
	CreatedAt        time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	ProcessingStatus ProcessingStatus       `bson:"status" json:"status"`
	Metadata         map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// NewMessage creates a new message with default values
func NewMessage(chatID primitive.ObjectID, content string, role MessageRole, msgType MessageType) *Message {
	now := time.Now()
	return &Message{
		ID:               primitive.NewObjectID(),
		ChatID:           chatID,
		Content:          content,
		Role:             role,
		Type:             msgType,
		CreatedAt:        now,
		UpdatedAt:        now,
		ProcessingStatus: StatusPending,
		Metadata:         make(map[string]interface{}),
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

// GetMetadataString retrieves a metadata value as string
func (m *Message) GetMetadataString(key string) string {
	if val, ok := m.GetMetadata(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetMetadataInt retrieves a metadata value as int
func (m *Message) GetMetadataInt(key string) int {
	if val, ok := m.GetMetadata(key); ok {
		// Try to convert from various numeric types
		switch v := val.(type) {
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float32:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

// GetMetadataBool retrieves a metadata value as bool
func (m *Message) GetMetadataBool(key string) bool {
	if val, ok := m.GetMetadata(key); ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// MarkAsSending marks the message as being sent to AI
func (m *Message) MarkAsSending() {
	m.ProcessingStatus = StatusSending
	m.UpdatedAt = time.Now()
}

// MarkAsComplete marks the message as completed processing
func (m *Message) MarkAsComplete() {
	m.ProcessingStatus = StatusComplete
	m.UpdatedAt = time.Now()
}

// MarkAsError marks the message as having an error
func (m *Message) MarkAsError(errorDetail string) {
	m.ProcessingStatus = StatusError
	m.UpdatedAt = time.Now()
	if errorDetail != "" {
		m.SetMetadata(MetaErrorDetail, errorDetail)
	}
}

// IsComplete checks if message processing is complete
func (m *Message) IsComplete() bool {
	return m.ProcessingStatus == StatusComplete
}
