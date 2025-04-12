# API Documentation

This document provides detailed information about the API endpoints, request/response formats, and SSE connection handling for the Go SSE AI Chat project.

## API Overview

The API follows RESTful principles with JSON for request and response payloads. Server-Sent Events (SSE) are used for real-time communication.

**Base URL**: `http://localhost:8080` (development) or your production domain

**API Version**: v1

## Authentication

Currently, the API uses a simple API key authentication method. Include the API key in the request header:

```
X-API-Key: your-api-key
```

## Common Headers

| Header | Description |
|--------|-------------|
| Content-Type | application/json |
| Accept | application/json |
| X-API-Key | API authentication key |

## Common Response Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request format or parameters |
| 401 | Unauthorized - Authentication failed |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server-side error |

## Endpoints

### Health Check

```
GET /health
```

Returns the status of the server and its dependencies.

**Response:**

```json
{
  "status": "ok",
  "timestamp": "2025-03-27T10:30:45Z",
  "services": {
    "mongodb": "connected",
    "ai_provider": "connected"
  },
  "version": "1.0.0"
}
```

### Chat Sessions

#### Create a new chat session

```
POST /api/v1/chats
```

Creates a new chat session.

**Request:**

```json
{
  "title": "Optional title for the chat"
}
```

**Response:**

```json
{
  "id": "chat_123456789",
  "title": "Optional title for the chat",
  "created_at": "2025-03-27T10:32:15Z"
}
```

#### List chat sessions

```
GET /api/v1/chats
```

Returns a list of all chat sessions.

**Response:**

```json
{
  "chats": [
    {
      "id": "chat_123456789",
      "title": "First conversation",
      "created_at": "2025-03-27T10:32:15Z",
      "last_message_at": "2025-03-27T10:45:30Z"
    },
    {
      "id": "chat_987654321",
      "title": "Another conversation",
      "created_at": "2025-03-27T11:12:22Z",
      "last_message_at": "2025-03-27T11:15:45Z"
    }
  ],
  "count": 2
}
```

#### Get a specific chat

```
GET /api/v1/chats/{chat_id}
```

Returns details of a specific chat session.

**Response:**

```json
{
  "id": "chat_123456789",
  "title": "First conversation",
  "created_at": "2025-03-27T10:32:15Z",
  "last_message_at": "2025-03-27T10:45:30Z",
  "message_count": 5
}
```

#### Delete a chat

```
DELETE /api/v1/chats/{chat_id}
```

Deletes a chat session and all associated messages.

**Response:**

```json
{
  "success": true,
  "message": "Chat deleted successfully"
}
```

### Messages

#### Send a message

```
POST /api/v1/chats/{chat_id}/messages
```

Sends a new message in a specific chat session.

**Request:**

```json
{
  "content": "Hello, how can you help me today?",
  "type": "text"
}
```

**Response:**

```json
{
  "id": "msg_123456789",
  "chat_id": "chat_123456789",
  "content": "Hello, how can you help me today?",
  "type": "text",
  "role": "user",
  "created_at": "2025-03-27T10:45:30Z"
}
```

#### Get messages from a chat

```
GET /api/v1/chats/{chat_id}/messages
```

Returns messages from a specific chat session.

**Query Parameters:**

| Parameter | Description | Default |
|-----------|-------------|---------|
| limit | Maximum number of messages to return | 50 |
| before | Return messages before this ID | - |

**Response:**

```json
{
  "messages": [
    {
      "id": "msg_123456789",
      "chat_id": "chat_123456789",
      "content": "Hello, how can you help me today?",
      "type": "text",
      "role": "user",
      "created_at": "2025-03-27T10:45:30Z"
    },
    {
      "id": "msg_987654321",
      "chat_id": "chat_123456789",
      "content": "I'm here to help you with any questions or tasks you have. What would you like assistance with today?",
      "type": "text",
      "role": "assistant",
      "created_at": "2025-03-27T10:45:35Z"
    }
  ],
  "has_more": false
}
```

## Server-Sent Events (SSE)

### Establishing an SSE Connection

```
GET /api/v1/chats/{chat_id}/stream
```

Establishes a Server-Sent Events (SSE) connection for receiving real-time messages.

**Headers:**

```
Accept: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

### SSE Event Format

Events are sent in the standard SSE format:

```
event: message
data: {"id":"msg_987654321","chat_id":"chat_123456789","content":"I'm here to help","type":"text","role":"assistant","created_at":"2025-03-27T10:45:35Z"}

```

### Event Types

| Event Type | Description |
|------------|-------------|
| message | New message or message chunk in the chat |
| error | Error message |
| ping | Keepalive message to maintain the connection |
| complete | Indicates that a streaming response is complete |

### Handling Stream Responses

When the AI is generating a response, chunks will be streamed as they become available:

```
event: message
data: {"id":"msg_abcdef123","chat_id":"chat_123456789","content":"I","type":"text","role":"assistant","created_at":"2025-03-27T10:46:00Z","is_chunk":true}

event: message
data: {"id":"msg_abcdef123","chat_id":"chat_123456789","content":"'m ","type":"text","role":"assistant","created_at":"2025-03-27T10:46:00Z","is_chunk":true}

event: message
data: {"id":"msg_abcdef123","chat_id":"chat_123456789","content":"analyzing ","type":"text","role":"assistant","created_at":"2025-03-27T10:46:00Z","is_chunk":true}

event: complete
data: {"id":"msg_abcdef123","chat_id":"chat_123456789"}
```

## Client Implementation

### JavaScript EventSource Example

```javascript
const chatId = "chat_123456789";
const eventSource = new EventSource(`/api/v1/chats/${chatId}/stream`);

// Handle new messages
eventSource.addEventListener('message', (event) => {
  const data = JSON.parse(event.data);
  console.log('Received message:', data);
  // Update UI with new message or chunk
});

// Handle errors
eventSource.addEventListener('error', (error) => {
  console.error('SSE Error:', error);
  // Implement reconnection logic
});

// Handle stream completion
eventSource.addEventListener('complete', (event) => {
  const data = JSON.parse(event.data);
  console.log('Stream completed for message:', data.id);
  // Update UI to indicate message completion
});

// Close the connection
function closeConnection() {
  eventSource.close();
}
```

### Reconnection Handling

For robust SSE connections, implement reconnection with exponential backoff:

```javascript
function connectSSE() {
  const chatId = "chat_123456789";
  const eventSource = new EventSource(`/api/v1/chats/${chatId}/stream`);
  
  // Set up event listeners as shown above
  
  eventSource.addEventListener('error', (error) => {
    console.error('SSE Error:', error);
    eventSource.close();
    
    // Implement exponential backoff
    const reconnectTime = calculateBackoff(retryCount);
    retryCount++;
    
    setTimeout(connectSSE, reconnectTime);
  });
  
  return eventSource;
}

function calculateBackoff(retry) {
  // Exponential backoff with jitter
  const baseDelay = 1000; // 1 second
  const maxDelay = 30000; // 30 seconds
  const delay = Math.min(maxDelay, baseDelay * Math.pow(2, retry));
  // Add jitter to prevent synchronized reconnections
  return delay + Math.random() * 1000;
}
```

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "bad_request",
    "message": "Invalid request parameters",
    "details": "Field 'content' is required"
  },
  "request_id": "req_123456789"
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| bad_request | The request was malformed or contained invalid parameters |
| unauthorized | Authentication failed |
| not_found | The requested resource was not found |
| rate_limited | Too many requests, try again later |
| ai_provider_error | Error from the AI provider |
| internal_error | Server-side error |

## Rate Limiting

API requests are rate-limited to prevent abuse. When rate limited, you'll receive a 429 status code:

```json
{
  "error": {
    "code": "rate_limited",
    "message": "Rate limit exceeded",
    "details": "Try again in 30 seconds"
  },
  "retry_after": 30
}
```

The `retry_after` field indicates how many seconds to wait before retrying.
