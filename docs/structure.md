# Project Structure

This document outlines the recommended folder structure and architectural organization for the Go SSE AI Chat project.

## Directory Structure

```
go-sse-ai-chat/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/               # Private application code
│   ├── api/                # API handlers
│   │   ├── chat.go         # Chat endpoint handlers
│   │   ├── health.go       # Health check endpoint
│   │   ├── middleware/     # Middleware components
│   │   └── routes.go       # API route definitions
│   ├── config/             # Configuration handling
│   │   └── config.go       # Configuration loader
│   ├── db/                 # Database layer
│   │   ├── models/         # Database models
│   │   └── repository/     # Data access repositories
│   ├── ai/                 # AI service integration
│   │   ├── openai.go       # OpenAI API client
│   │   ├── anthropic.go    # Anthropic API client
│   │   └── service.go      # AI service interface
│   ├── sse/                # SSE implementation
│   │   ├── client.go       # SSE client connection
│   │   ├── broker.go       # SSE message broker
│   │   └── handler.go      # SSE request handler
│   └── domain/             # Core domain logic
│       ├── chat.go         # Chat domain logic
│       └── message.go      # Message domain logic
├── pkg/                    # Public library code
│   ├── logger/             # Logging utilities
│   └── utils/              # Shared utilities
├── static/                 # Static frontend assets
│   ├── css/
│   ├── js/
│   └── index.html
├── tests/                  # Test files
│   ├── integration/
│   └── unit/
├── docs/                   # Documentation
│   ├── issues.md
│   ├── structure.md
│   ├── api.md
│   └── deployment.md
├── scripts/                # Build and deployment scripts
├── .env.example            # Example environment variables
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project overview and instructions
```

## Architectural Layers

The application follows a clean architecture pattern with distinct layers:

1. **Presentation Layer** (`api/`) - Handles HTTP requests and API endpoints
2. **Domain Layer** (`domain/`) - Contains the core business logic
3. **Data Access Layer** (`db/repository/`) - Handles database operations
4. **Infrastructure Layer** (`sse/`, `ai/`, `config/`) - Provides technical capabilities 

## Key Components

### HTTP Server (Gin)

The application uses Gin as the HTTP server framework for handling routes and middleware.

```go
// Example route setup with Gin
func SetupRouter(config Config) *gin.Engine {
    router := gin.Default()
    
    router.GET("/health", api.HealthCheck)
    router.POST("/api/chat", api.HandleChatMessage)
    router.GET("/api/events", sse.HandleSSEConnection)
    
    return router
}
```

### SSE Implementation

The SSE implementation consists of several core components:

1. **Client** - Represents a connected client with a unique ID
2. **Broker** - Manages connected clients and dispatches messages
3. **Handler** - HTTP handler that establishes and maintains SSE connections

```go
// Example SSE client structure
type Client struct {
    ID       string
    Messages chan []byte
}

// Example broker structure
type Broker struct {
    clients    map[string]*Client
    register   chan *Client
    unregister chan *Client
    messages   chan *Message
}
```

### Database Layer

MongoDB is used for data persistence with a repository pattern:

```go
// Example repository interface
type ChatRepository interface {
    CreateChat(ctx context.Context, chat *models.Chat) (string, error)
    GetChatByID(ctx context.Context, id string) (*models.Chat, error)
    SaveMessage(ctx context.Context, chatID string, message *models.Message) error
    GetMessagesByChatID(ctx context.Context, chatID string) ([]*models.Message, error)
}
```

### AI Service Integration

The AI integration is designed with an interface to support multiple providers:

```go
// Example AI service interface
type AIService interface {
    GenerateResponse(ctx context.Context, prompt string, stream bool) (chan string, error)
    CompleteStream(ctx context.Context, messages []Message) (chan string, error)
}
```

## Data Flow

1. Client sends a message via HTTP POST to `/api/chat`
2. Server processes the request and stores the message in MongoDB
3. Server sends the message to the AI service for processing
4. AI service streams the response back to the server
5. The SSE broker sends message chunks to the connected client in real-time
6. Client renders the response as it arrives

## Configuration

The application uses environment variables for configuration, which can be set via a `.env` file, environment variables, or command-line flags:

```
# Example environment variables
SERVER_PORT=8080
MONGODB_URI=mongodb://localhost:27017/sse-chat
AI_PROVIDER=openai
OPENAI_API_KEY=your-api-key
```

## Deployment Options

The application can be deployed in several ways:

1. **Docker** - Using the provided Dockerfile and docker-compose.yml
2. **Kubernetes** - For scalable deployments (configuration in k8s/ directory)
3. **Traditional deployment** - Running the binary directly on a server
