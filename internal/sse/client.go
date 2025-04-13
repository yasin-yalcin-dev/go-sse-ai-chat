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
	"fmt"
	"net/http"
	"time"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
)

// Client represents a connected SSE client
type Client struct {
	ID           string
	Connection   http.ResponseWriter
	MessageChan  chan []byte
	ConnectedAt  time.Time
	LastActivity time.Time
	IsClosed     bool
	Ctx          context.Context
	Cancel       context.CancelFunc
	Broker       *Broker
}

// NewClient creates a new SSE client
func NewClient(id string, w http.ResponseWriter, broker *Broker) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		ID:           id,
		Connection:   w,
		MessageChan:  make(chan []byte, 256), // Buffer for messages
		ConnectedAt:  time.Now(),
		LastActivity: time.Now(),
		IsClosed:     false,
		Ctx:          ctx,
		Cancel:       cancel,
		Broker:       broker,
	}
}

// Send sends a message to the client
func (c *Client) Send(message []byte) error {
	// Check if client is already closed
	if c.IsClosed {
		return fmt.Errorf("client %s is closed", c.ID)
	}

	// Try to send or timeout
	select {
	case c.MessageChan <- message:
		c.LastActivity = time.Now()
		return nil
	case <-time.After(5 * time.Second): // Timeout if channel is full
		return fmt.Errorf("send timeout for client %s", c.ID)
	}
}

// Listen starts listening for messages and sends them to the client
func (c *Client) Listen() {
	// Set necessary headers for SSE
	c.setupHeaders()

	// Create a flush for the client
	flusher, ok := c.Connection.(http.Flusher)
	if !ok {
		logger.Errorf("Could not initialize SSE connection: %s - client doesn't support flushing", c.ID)
		c.Close()
		return
	}

	// Signal to the broker that this client is ready
	c.Broker.Register <- c

	// Send an initial ping to establish connection
	if err := c.sendPing(); err != nil {
		logger.Errorf("Failed to send initial ping to client %s: %v", c.ID, err)
		c.Close()
		return
	}

	// Create keepalive ticker
	keepalive := time.NewTicker(c.Broker.KeepaliveInterval)
	defer keepalive.Stop()

	// Listen for messages and send them to the client
	for {
		select {
		case <-c.Ctx.Done():
			// Context was canceled, exit
			logger.Debugf("Context canceled for client %s", c.ID)
			c.Close()
			return

		case <-keepalive.C:
			// Send keepalive ping
			if err := c.sendPing(); err != nil {
				logger.Warnf("Failed to send keepalive to client %s: %v", c.ID, err)
				c.Close()
				return
			}

		case msg, ok := <-c.MessageChan:
			if !ok {
				// Channel was closed
				logger.Debugf("Message channel closed for client %s", c.ID)
				c.Close()
				return
			}

			// Write message to the connection
			if _, err := fmt.Fprintf(c.Connection, "data: %s\n\n", msg); err != nil {
				logger.Warnf("Failed to send message to client %s: %v", c.ID, err)
				c.Close()
				return
			}

			// Flush to send data immediately
			flusher.Flush()
			c.LastActivity = time.Now()
		}
	}
}

// Close closes the client connection
func (c *Client) Close() {
	// Prevent double closing
	if c.IsClosed {
		return
	}

	c.IsClosed = true
	c.Cancel() // Cancel the context

	// Tell broker to unregister this client
	c.Broker.Unregister <- c

	// Close message channel
	close(c.MessageChan)

	logger.Infof("SSE client disconnected: %s (connected for %v)",
		c.ID, time.Since(c.ConnectedAt))
}

// setupHeaders sets the necessary headers for SSE
func (c *Client) setupHeaders() {
	header := c.Connection.Header()

	// Required SSE headers
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no") // Disable buffering in Nginx

	// CORS headers (if needed)
	header.Set("Access-Control-Allow-Origin", "*")
}

// sendEvent sends a named event to the client
func (c *Client) sendEvent(event string, data []byte) error {
	flusher, ok := c.Connection.(http.Flusher)
	if !ok {
		return fmt.Errorf("client doesn't support flushing")
	}

	if _, err := fmt.Fprintf(c.Connection, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}

	flusher.Flush()
	c.LastActivity = time.Now()
	return nil
}

// sendPing sends a keepalive ping
func (c *Client) sendPing() error {
	timestamp := time.Now().Unix()
	ping := fmt.Sprintf(`{"time":%d}`, timestamp)
	return c.sendEvent("ping", []byte(ping))
}
