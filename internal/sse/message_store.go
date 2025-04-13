/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package sse

import (
	"sync"
	"time"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// MessageStore keeps track of recent messages for replay on reconnect
type MessageStore struct {
	// Maps chat ID to a list of recent messages for that chat
	messages map[string][]*StoredMessage

	// Maximum number of messages to store per chat
	maxMessagesPerChat int

	// Maximum age of messages to keep (older messages are pruned)
	maxMessageAge time.Duration

	// Last cleanup time
	lastCleanup time.Time

	// Mutex for thread safety
	mutex sync.RWMutex
}

// StoredMessage represents a message stored for potential replay
type StoredMessage struct {
	ID        string    // Unique message ID
	Event     string    // Event type
	Data      []byte    // Message data
	Timestamp time.Time // When the message was sent
}

// NewMessageStore creates a new message store
func NewMessageStore(maxMessagesPerChat int, maxMessageAge time.Duration) *MessageStore {
	return &MessageStore{
		messages:           make(map[string][]*StoredMessage),
		maxMessagesPerChat: maxMessagesPerChat,
		maxMessageAge:      maxMessageAge,
		lastCleanup:        time.Now(),
	}
}

// StoreMessage adds a message to the store
func (s *MessageStore) StoreMessage(chatID string, messageID string, event string, data []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create the chat's message queue if it doesn't exist
	if _, exists := s.messages[chatID]; !exists {
		s.messages[chatID] = make([]*StoredMessage, 0, s.maxMessagesPerChat)
	}

	// Add the new message
	s.messages[chatID] = append(s.messages[chatID], &StoredMessage{
		ID:        messageID,
		Event:     event,
		Data:      data,
		Timestamp: time.Now(),
	})

	// Trim if we have too many messages
	if len(s.messages[chatID]) > s.maxMessagesPerChat {
		// Remove oldest message
		s.messages[chatID] = s.messages[chatID][1:]
	}

	// Periodically clean up old messages (every 5 minutes)
	if time.Since(s.lastCleanup) > 5*time.Minute {
		go s.cleanup()
	}
}

// GetRecentMessages retrieves recent messages for a chat
func (s *MessageStore) GetRecentMessages(chatID string, since time.Time) []*StoredMessage {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	messages, exists := s.messages[chatID]
	if !exists {
		return []*StoredMessage{}
	}

	// Find messages newer than the given time
	var recentMessages []*StoredMessage
	for _, msg := range messages {
		if msg.Timestamp.After(since) {
			recentMessages = append(recentMessages, msg)
		}
	}

	return recentMessages
}

// cleanup removes old messages
func (s *MessageStore) cleanup() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.lastCleanup = time.Now()
	cutoff := time.Now().Add(-s.maxMessageAge)

	for chatID, messages := range s.messages {
		var newMessages []*StoredMessage

		for _, msg := range messages {
			if msg.Timestamp.After(cutoff) {
				newMessages = append(newMessages, msg)
			}
		}

		// Update or delete the chat's message list
		if len(newMessages) > 0 {
			s.messages[chatID] = newMessages
		} else {
			delete(s.messages, chatID)
		}
	}

	logger.Debugf("Cleaned up message store, now tracking %d chats", len(s.messages))
}
