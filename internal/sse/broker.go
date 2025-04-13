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
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// Broker manages SSE clients and message distribution
type Broker struct {
	// Client management
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client

	// Message broadcasting
	Broadcast chan *Message

	// Message store for reconnection replay
	messageStore *MessageStore

	// Configuration
	MaxClients        int
	KeepaliveInterval time.Duration
	MaxRetryAttempts  int
	RetryDelay        time.Duration

	// Synchronization
	mutex sync.RWMutex
}

// Message represents a message to be sent to clients
type Message struct {
	ID       string          // Unique message ID
	ChatID   string          // Chat ID this message belongs to
	Event    string          // Event name
	Data     json.RawMessage // Message data as JSON
	Target   string          // Target client ID (empty for broadcast to whole chat)
	Attempts int             // Number of delivery attempts made
}

// NewBroker creates a new SSE broker
func NewBroker(maxClients int, keepaliveInterval time.Duration) *Broker {
	// Create message store that keeps messages for 5 minutes, up to 50 per chat
	messageStore := NewMessageStore(50, 5*time.Minute)

	return &Broker{
		Clients:           make(map[string]*Client),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		Broadcast:         make(chan *Message, 256), // Buffer for messages
		messageStore:      messageStore,
		MaxClients:        maxClients,
		KeepaliveInterval: keepaliveInterval,
		MaxRetryAttempts:  3,                      // Retry failed messages 3 times
		RetryDelay:        500 * time.Millisecond, // Wait 500ms between retries
	}
}

// Start starts the broker
func (b *Broker) Start(ctx context.Context) {
	logger.Info("Starting SSE broker")

	for {
		select {
		case <-ctx.Done():
			// Application is shutting down
			logger.Info("Shutting down SSE broker")
			b.closeAllClients()
			return

		case client := <-b.Register:
			// New client connected
			b.registerClient(client)

		case client := <-b.Unregister:
			// Client disconnected
			b.unregisterClient(client)

		case message := <-b.Broadcast:
			// New message to broadcast
			b.processMessage(message)
		}
	}
}

// registerClient registers a new client
func (b *Broker) registerClient(client *Client) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check max clients limit
	if len(b.Clients) >= b.MaxClients {
		logger.Warnf("Max SSE clients reached (%d), rejecting new connection", b.MaxClients)
		client.Close()
		return
	}

	// Get the chat ID from the client ID (format: chatID_clientUUID)
	chatID := getChatIDFromClientID(client.ID)

	// Add client to the map
	b.Clients[client.ID] = client
	logger.Infof("SSE client connected: %s (total clients: %d)", client.ID, len(b.Clients))

	// Check if we need to replay messages for this chat
	if chatID != "" {
		// Send last 5 minutes of messages
		b.replayMessages(client, chatID, time.Now().Add(-5*time.Minute))
	}
}

// unregisterClient removes a client
func (b *Broker) unregisterClient(client *Client) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if client exists
	if _, exists := b.Clients[client.ID]; !exists {
		return
	}

	// Remove client from the map
	delete(b.Clients, client.ID)
	logger.Debugf("Unregistered SSE client: %s (remaining clients: %d)", client.ID, len(b.Clients))
}

// processMessage handles a new message, storing it and delivering it
func (b *Broker) processMessage(message *Message) {
	// Extract chat ID from message or target
	chatID := message.ChatID
	if chatID == "" && message.Target != "" {
		chatID = getChatIDFromClientID(message.Target)
	}

	// Store the message for replay if it has a chat ID and message ID
	if chatID != "" && message.ID != "" {
		// Convert message to JSON for storage
		messageJSON, err := json.Marshal(message.Data)
		if err == nil {
			b.messageStore.StoreMessage(chatID, message.ID, message.Event, messageJSON)
		}
	}

	// Deliver the message
	b.deliverMessage(message)
}

// deliverMessage sends a message to the targeted client(s)
func (b *Broker) deliverMessage(message *Message) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Convert message to JSON
	messageJSON, err := json.Marshal(message.Data)
	if err != nil {
		logger.Errorf("Failed to marshal message: %v", err)
		return
	}

	// Check if message is targeted to a specific client
	if message.Target != "" {
		// Send to specific client
		client, exists := b.Clients[message.Target]
		if exists {
			if err := client.Send(messageJSON); err != nil {
				logger.Warnf("Failed to send message to client %s: %v", client.ID, err)
				// If sending failed and we haven't reached max retries, queue for retry
				if message.Attempts < b.MaxRetryAttempts {
					go b.retryMessage(message)
				}
			}
		}
	} else if message.ChatID != "" {
		// Send to all clients in this chat
		for clientID, client := range b.Clients {
			// Check if this client belongs to the target chat
			if getChatIDFromClientID(clientID) == message.ChatID {
				if err := client.Send(messageJSON); err != nil {
					logger.Warnf("Failed to send message to client %s: %v", client.ID, err)
				}
			}
		}
	} else {
		// Broadcast to all clients
		for _, client := range b.Clients {
			if err := client.Send(messageJSON); err != nil {
				logger.Warnf("Failed to send message to client %s: %v", client.ID, err)
			}
		}
	}
}

// retryMessage attempts to resend a failed message after a delay
func (b *Broker) retryMessage(message *Message) {
	// Increment attempt count
	message.Attempts++

	// Wait before retrying
	time.Sleep(b.RetryDelay)

	// Try to send again
	logger.Debugf("Retrying message delivery (attempt %d/%d)",
		message.Attempts, b.MaxRetryAttempts)

	b.Broadcast <- message
}

// replayMessages sends recent messages to a newly connected client
func (b *Broker) replayMessages(client *Client, chatID string, since time.Time) {
	// Get recent messages for this chat
	messages := b.messageStore.GetRecentMessages(chatID, since)

	if len(messages) > 0 {
		logger.Infof("Replaying %d messages for client %s", len(messages), client.ID)

		// Send a notification that we're replaying messages
		replayStart := map[string]interface{}{
			"type":  "replay_start",
			"count": len(messages),
		}
		replayStartJSON, _ := json.Marshal(replayStart)
		client.sendEvent("control", replayStartJSON)

		// Send each message
		for _, msg := range messages {
			client.sendEvent(msg.Event, msg.Data)
		}

		// Send replay complete notification
		replayEnd := map[string]interface{}{
			"type": "replay_end",
		}
		replayEndJSON, _ := json.Marshal(replayEnd)
		client.sendEvent("control", replayEndJSON)
	}
}

// closeAllClients closes all client connections
func (b *Broker) closeAllClients() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	logger.Infof("Closing all SSE connections (%d clients)", len(b.Clients))

	for _, client := range b.Clients {
		client.Close()
	}

	// Clear the clients map
	b.Clients = make(map[string]*Client)
}

// GetClientCount returns the number of connected clients
func (b *Broker) GetClientCount() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return len(b.Clients)
}

// GetClientsInChat returns the number of clients connected to a specific chat
func (b *Broker) GetClientsInChat(chatID string) int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	count := 0
	for clientID := range b.Clients {
		if getChatIDFromClientID(clientID) == chatID {
			count++
		}
	}

	return count
}

// SendToClient sends a message to a specific client
func (b *Broker) SendToClient(clientID string, messageID string, event string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	message := &Message{
		ID:     messageID,
		Event:  event,
		Data:   dataJSON,
		Target: clientID,
	}

	// If this is a targeted message, try to extract the chat ID from client ID
	message.ChatID = getChatIDFromClientID(clientID)

	b.Broadcast <- message
	return nil
}

// SendToChat sends a message to all clients in a chat
func (b *Broker) SendToChat(chatID string, messageID string, event string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	b.Broadcast <- &Message{
		ID:     messageID,
		ChatID: chatID,
		Event:  event,
		Data:   dataJSON,
	}

	return nil
}

// BroadcastToAll sends a message to all clients
func (b *Broker) BroadcastToAll(event string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	b.Broadcast <- &Message{
		Event: event,
		Data:  dataJSON,
	}

	return nil
}

// Helper function to extract chat ID from client ID
func getChatIDFromClientID(clientID string) string {
	parts := strings.Split(clientID, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
