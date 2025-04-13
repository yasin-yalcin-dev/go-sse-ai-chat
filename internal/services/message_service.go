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
	"errors"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/ai"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/repository"
	customerrors "github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/errors"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageServiceImpl implements the MessageService interface
type MessageServiceImpl struct {
	messageRepo    repository.MessageRepository
	chatRepo       repository.ChatRepository
	contextManager *ai.ContextManager
}

// NewMessageService creates a new message service
func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	// Create context manager with default values
	contextManager := ai.NewContextManager(
		ai.DefaultMaxContextMessages,
		"You are a helpful assistant.", // Default system prompt
	)
	return &MessageServiceImpl{
		messageRepo:    messageRepo,
		chatRepo:       chatRepo,
		contextManager: contextManager,
	}
}

// CreateMessage creates a new message
func (s *MessageServiceImpl) CreateMessage(ctx context.Context, chatID string, content string, role models.MessageRole, msgType models.MessageType) (*models.Message, error) {
	// Validate and sanitize message content
	sanitizedContent, err := utils.ValidateAndSanitizeMessage(content)
	if err != nil {
		return nil, customerrors.NewValidationError(err.Error(), err)
	}

	chatObjID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, customerrors.NewBadRequestError("Invalid chat ID format", err)
	}

	// Verify the chat exists
	chat, err := s.chatRepo.FindByID(ctx, chatObjID)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, customerrors.NewNotFoundError("Chat not found", nil)
	}

	// Create the message with sanitized content
	message := models.NewMessage(chatObjID, sanitizedContent, role, msgType)

	// Add creation metadata
	message.SetMetadata(models.MetaClientIP, getClientIP(ctx))
	message.SetMetadata(models.MetaUserAgent, getUserAgent(ctx))

	// Add code language metadata if it's a code message
	if msgType == models.TypeCode {
		if lang, ok := ctx.Value("code_language").(string); ok && lang != "" {
			message.SetMetadata(models.MetaCodeLanguage, lang)
		}
	}

	// Add reply metadata if replying to a message
	if replyToID, ok := ctx.Value("reply_to").(string); ok && replyToID != "" {
		message.SetMetadata(models.MetaReplyTo, replyToID)
	}

	// For assistant messages, provide initial pending status
	if role == models.RoleAssistant {
		message.MarkAsSending()
	} else {
		// User messages are automatically complete
		message.MarkAsComplete()
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Update the chat's message count
	if err := s.chatRepo.IncrementMessageCount(ctx, chatObjID); err != nil {
		// Log the error but don't fail the operation
		logger.Errorf("Failed to increment message count: %v", err)
	}

	return message, nil
}

// Helper functions to get request metadata
func getClientIP(ctx context.Context) string {
	// Try to get client IP from context if available
	// This would require passing it from the handler layer
	if ip, ok := ctx.Value("client_ip").(string); ok {
		return ip
	}
	return "unknown"
}

func getUserAgent(ctx context.Context) string {
	// Try to get user agent from context if available
	if ua, ok := ctx.Value("user_agent").(string); ok {
		return ua
	}
	return "unknown"
}

// GetMessageByID retrieves a message by its ID
func (s *MessageServiceImpl) GetMessageByID(ctx context.Context, id string) (*models.Message, error) {
	msgID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	message, err := s.messageRepo.FindByID(ctx, msgID)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.New("message not found")
	}

	return message, nil
}

// GetChatMessages retrieves paginated messages for a chat
func (s *MessageServiceImpl) GetChatMessages(ctx context.Context, chatID string, page, pageSize int) ([]*models.Message, int64, error) {
	chatObjID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, 0, err
	}

	// Verify the chat exists
	chat, err := s.chatRepo.FindByID(ctx, chatObjID)
	if err != nil {
		return nil, 0, err
	}

	if chat == nil {
		return nil, 0, errors.New("chat not found")
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	messages, err := s.messageRepo.FindByChatID(ctx, chatObjID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.messageRepo.CountByChatID(ctx, chatObjID)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// DeleteMessage deletes a message
func (s *MessageServiceImpl) DeleteMessage(ctx context.Context, id string) error {
	msgID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Get the message to check which chat it belongs to
	message, err := s.messageRepo.FindByID(ctx, msgID)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.New("message not found")
	}

	// Delete the message
	if err := s.messageRepo.Delete(ctx, msgID); err != nil {
		return err
	}

	// Update the chat's message count (decrement)
	// This could be implemented with a new repository method
	// For now, we'll just note that it should be done
	// TODO: Implement chat.DecrementMessageCount

	return nil
}

// GetMessageContext retrieves AI-ready message context for a chat
func (s *MessageServiceImpl) GetMessageContext(ctx context.Context, chatID string, maxMessages int) ([]ai.AIMessage, error) {
	chatObjID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, customerrors.NewBadRequestError("Invalid chat ID format", err)
	}

	// Get chat to verify it exists
	chat, err := s.chatRepo.FindByID(ctx, chatObjID)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, customerrors.NewNotFoundError("Chat not found", nil)
	}

	// Adjust max messages if needed
	if maxMessages <= 0 {
		maxMessages = s.contextManager.MaxContextMessages
	}

	// Get messages for this chat
	messages, err := s.messageRepo.FindByChatID(ctx, chatObjID, maxMessages, 0)
	if err != nil {
		return nil, err
	}

	// Create AI context from messages
	aiContext := s.contextManager.CreateContext(messages)

	return aiContext, nil
}
