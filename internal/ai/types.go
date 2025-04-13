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
	"context"
	"time"
)

// ProviderType represents the type of AI provider
type ProviderType string

const (
	// ProviderOpenAI is OpenAI provider (e.g. GPT models)
	ProviderOpenAI ProviderType = "openai"

	// ProviderAnthropic is Anthropic provider (e.g. Claude models)
	ProviderAnthropic ProviderType = "anthropic"
)

// AIMessage represents a message in the format required by AI providers
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest contains parameters for an AI completion request
type CompletionRequest struct {
	Messages    []AIMessage `json:"messages"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature float64     `json:"temperature,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
}

// CompletionResponse represents a response from an AI provider
type CompletionResponse struct {
	Content      string    `json:"content"`
	Created      time.Time `json:"created"`
	Model        string    `json:"model"`
	IsComplete   bool      `json:"is_complete"`
	FinishReason string    `json:"finish_reason,omitempty"`
}

// Provider defines the interface for AI providers
type Provider interface {
	// GetType returns the type of AI provider
	GetType() ProviderType

	// GetModelName returns the name of the model being used
	GetModelName() string

	// Complete generates a completion for the given request
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	// CompleteStream generates a streaming completion for the given request
	CompleteStream(ctx context.Context, req CompletionRequest) (<-chan *CompletionResponse, error)
}
