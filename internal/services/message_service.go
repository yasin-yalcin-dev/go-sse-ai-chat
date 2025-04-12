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

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/models"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageServiceImpl implements the MessageService interface
type MessageServiceImpl struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
}

// NewMessageService creates a new message service
func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) MessageService {
	return &MessageServiceImpl{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
	}
}

// CreateMessage creates a new message
func (s *MessageServiceImpl) CreateMessage(ctx context.Context, chatID string, content string, role models.MessageRole, msgType models.MessageType) (*models.Message, error) {
	chatObjID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, err
	}

	// Verify the chat exists
	chat, err := s.chatRepo.FindByID(ctx, chatObjID)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("chat not found")
	}

	// Create the message
	message := models.NewMessage(chatObjID, content, role, msgType)

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Update the chat's message count
	if err := s.chatRepo.IncrementMessageCount(ctx, chatObjID); err != nil {
		// Log the error but don't fail the operation
		// TODO: Add proper logging
		// logger.Errorf("Failed to increment message count: %v", err)
	}

	return message, nil
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
