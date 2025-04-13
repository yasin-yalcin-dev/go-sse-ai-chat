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

	// Configuration
	MaxClients        int
	KeepaliveInterval time.Duration

	// Synchronization
	mutex sync.RWMutex
}

// Message represents a message to be sent to clients
type Message struct {
	Event  string          // Event name
	Data   json.RawMessage // Message data as JSON
	Target string          // Target client ID (empty for broadcast)
}

// NewBroker creates a new SSE broker
func NewBroker(maxClients int, keepaliveInterval time.Duration) *Broker {
	return &Broker{
		Clients:           make(map[string]*Client),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		Broadcast:         make(chan *Message, 256), // Buffer for messages
		MaxClients:        maxClients,
		KeepaliveInterval: keepaliveInterval,
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
			b.broadcastMessage(message)
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

	// Add client to the map
	b.Clients[client.ID] = client
	logger.Infof("SSE client connected: %s (total clients: %d)", client.ID, len(b.Clients))
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

// broadcastMessage sends a message to clients
func (b *Broker) broadcastMessage(message *Message) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Convert message to JSON
	messageJSON, err := json.Marshal(message.Data)
	if err != nil {
		logger.Errorf("Failed to marshal message: %v", err)
		return
	}

	// Check if message is targeted or broadcast
	if message.Target != "" {
		// Send to specific client
		if client, exists := b.Clients[message.Target]; exists {
			if err := client.Send(messageJSON); err != nil {
				logger.Warnf("Failed to send message to client %s: %v", client.ID, err)
			}
		}
	} else {
		// Send to all clients
		for _, client := range b.Clients {
			if err := client.Send(messageJSON); err != nil {
				logger.Warnf("Failed to send message to client %s: %v", client.ID, err)
			}
		}
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

// SendToClient sends a message to a specific client
func (b *Broker) SendToClient(clientID string, event string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	b.Broadcast <- &Message{
		Event:  event,
		Data:   dataJSON,
		Target: clientID,
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
