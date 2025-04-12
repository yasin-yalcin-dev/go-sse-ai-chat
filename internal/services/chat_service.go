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

// ChatServiceImpl implements the ChatService interface
type ChatServiceImpl struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
}

// NewChatService creates a new chat service
func NewChatService(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository) ChatService {
	return &ChatServiceImpl{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

// CreateChat creates a new chat session
func (s *ChatServiceImpl) CreateChat(ctx context.Context, title string) (*models.Chat, error) {
	chat := models.NewChat(title)

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

// GetChatByID retrieves a chat by its ID
func (s *ChatServiceImpl) GetChatByID(ctx context.Context, id string) (*models.Chat, error) {
	chatID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("chat not found")
	}

	return chat, nil
}

// ListChats retrieves a paginated list of chats
func (s *ChatServiceImpl) ListChats(ctx context.Context, page, pageSize int) ([]*models.Chat, int64, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	chats, err := s.chatRepo.FindAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.chatRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

// UpdateChat updates a chat's title
func (s *ChatServiceImpl) UpdateChat(ctx context.Context, id string, title string) (*models.Chat, error) {
	chatID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	chat, err := s.chatRepo.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("chat not found")
	}

	chat.Title = title

	if err := s.chatRepo.Update(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

// DeleteChat deletes a chat and all its messages
func (s *ChatServiceImpl) DeleteChat(ctx context.Context, id string) error {
	chatID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Delete all messages first
	if err := s.messageRepo.DeleteByChatID(ctx, chatID); err != nil {
		return err
	}

	// Then delete the chat
	return s.chatRepo.Delete(ctx, chatID)
}
