# Go SSE AI Chat

A real-time AI-powered chat application using Server-Sent Events (SSE) for streaming responses.

## Overview

This project is a high-performance chat application that connects to AI providers (OpenAI/Anthropic) and streams responses in real-time to the client using Server-Sent Events. It's built with Go and follows clean architecture principles.

## Key Features

- **Real-time Streaming**: Uses Server-Sent Events (SSE) to stream AI responses as they're generated
- **Multiple AI Providers**: Supports both OpenAI and Anthropic APIs
- **Persistent Storage**: Stores conversations and messages in MongoDB
- **Clean Architecture**: Follows the repository pattern and dependency injection principles
- **High Performance**: Built with Go for efficient handling of concurrent connections
- **Responsive UI**: Modern, responsive web interface for seamless chatting experience

## Project Structure

```
go-sse-ai-chat/
├── cmd/              # Application entry points
│   └── api/          # API server
├── internal/         # Private application code
│   ├── config/       # Configuration management
│   ├── db/           # Database connections and migrations
│   ├── handlers/     # HTTP request handlers
│   ├── middleware/   # HTTP middleware components
│   ├── models/       # Data models
│   ├── repository/   # Data access layer
│   └── services/     # Business logic
├── pkg/              # Public packages for reuse
│   ├── errors/       # Error handling utilities
│   ├── logger/       # Logging utilities
│   └── utils/        # General utilities
├── static/           # Static assets for frontend
├── docs/             # Documentation
├── .env.example      # Example environment variables
└── README.md         # This file
```

## Getting Started

### Prerequisites

- Go 1.20+
- MongoDB 5.0+
- OpenAI API key or Anthropic API key

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-sse-ai-chat.git
   cd go-sse-ai-chat
   ```

2. Create and configure the environment file:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration details
   ```

3. Build and run the application:
   ```bash
   go build -o app ./cmd/api
   ./app
   ```

### Configuration

See `.env.example` for all available configuration options.

## Development

### Running Tests

```bash
go test ./...
```

### Development Mode

```bash
go run ./cmd/api/main.go
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.