/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package ai

import (
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// Default context parameters
const (
	DefaultMaxContextMessages = 20   // Maximum number of messages to include in context
	DefaultMaxTokens          = 4096 // Maximum number of tokens for the response
	DefaultTemperature        = 0.7  // Default temperature for response generation
)

// ContextManager handles creating and maintaining AI message context
type ContextManager struct {
	// Configuration
	MaxContextMessages int
	SystemPrompt       string
}

// NewContextManager creates a new context manager
func NewContextManager(maxMessages int, systemPrompt string) *ContextManager {
	if maxMessages <= 0 {
		maxMessages = DefaultMaxContextMessages
	}

	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant."
	}

	return &ContextManager{
		MaxContextMessages: maxMessages,
		SystemPrompt:       systemPrompt,
	}
}

// CreateContext creates an AI context from chat messages
func (cm *ContextManager) CreateContext(messages []*models.Message) []AIMessage {
	if len(messages) == 0 {
		// If no messages, just return system prompt
		return []AIMessage{
			{Role: "system", Content: cm.SystemPrompt},
		}
	}

	// Start with system message
	context := []AIMessage{
		{Role: "system", Content: cm.SystemPrompt},
	}

	// Limit number of messages to include in context
	startIdx := 0
	if len(messages) > cm.MaxContextMessages {
		startIdx = len(messages) - cm.MaxContextMessages
		logger.Debugf("Limiting context to %d of %d messages", cm.MaxContextMessages, len(messages))
	}

	// Add messages to context
	for i := startIdx; i < len(messages); i++ {
		msg := messages[i]

		// Map message role to AI role
		role := "user"
		if msg.Role == models.RoleAssistant {
			role = "assistant"
		} else if msg.Role == models.RoleSystem {
			role = "system"
		}

		context = append(context, AIMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	return context
}

// CreateCompletionRequest creates a completion request with the given messages
func (cm *ContextManager) CreateCompletionRequest(messages []*models.Message, stream bool) CompletionRequest {
	// Create AI message context
	context := cm.CreateContext(messages)

	// Create completion request
	return CompletionRequest{
		Messages:    context,
		MaxTokens:   DefaultMaxTokens,
		Temperature: DefaultTemperature,
		Stream:      stream,
	}
}
