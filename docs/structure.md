# Project Structure

This document outlines the recommended folder structure and architectural organization for the Go SSE AI Chat project.

## Directory Structure (Updated)

```
go-sse-ai-chat/
├── cmd/
│   └── api/
│       ├── main.go         # Application entry point
│       └── routes.go       # API route configuration
├── internal/               # Private application code
│   ├── config/             # Configuration handling
│   │   └── config.go       # Configuration loader
│   ├── db/                 # Database layer
│   │   └── mongodb/        # MongoDB implementation 
│   │       ├── connection.go  # Connection manager
│   │       ├── client.go      # DB client utilities
│   │       ├── options.go     # DB options
│   │       └── indexes.go     # Indexes management
│   ├── handlers/           # HTTP request handlers
│   │   ├── handler.go      # Base handler
│   │   ├── chat_handler.go # Chat endpoints
│   │   ├── message_handler.go # Message endpoints
│   │   └── system_handler.go  # System endpoints
│   ├── middleware/         # Middleware components
│   │   ├── cors.go         # CORS middleware
│   │   ├── logger.go       # Logging middleware
│   │   └── error.go        # Error handling middleware
│   ├── models/             # Data models
│   │   ├── chat.go         # Chat model
│   │   ├── message.go      # Message model
│   │   └── dto/            # Data Transfer Objects
│   │       ├── common.go   # Common DTOs
│   │       ├── chat_dto.go # Chat DTOs
│   │       └── message_dto.go # Message DTOs
│   ├── repository/         # Repository interfaces
│   │   ├── repository.go   # Repository definitions
│   │   └── mongodb/       # MongoDB implementations
│   │       ├── chat_repository.go    # Chat repo
│   │       └── message_repository.go # Message repo
│   └── services/           # Business logic
│       ├── chat_service.go  # Chat operations
│       └── message_service.go # Message operations
├── pkg/                    # Public library code
│   ├── errors/             # Error handling utilities
│   │   └── errors.go       # Custom error types
│   ├── logger/             # Logging utilities
│   │   └── logger.go       # Logger implementation
│   └── utils/              # Shared utilities
├── static/                 # Static frontend assets
│   ├── css/
│   ├── js/
│   └── index.html
├── docs/                   # Documentation
│   ├── issues.md           # Project issues
│   ├── progress.md         # Development progress
│   ├── structure.md        # Project structure
│   ├── api.md              # API documentation
│   └── setup.md            # Setup instructions
├── .env.example            # Example environment variables
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project overview and instructions
```

## Architectural Layers

The application follows a clean architecture pattern with distinct layers:

1. **Presentation Layer** (`handlers/`) - Handles HTTP requests and API endpoints
2. **Domain Layer** (`models/`, `services/`) - Contains business logic and data structures
3. **Data Access Layer** (`repository/`) - Handles database operations
4. **Infrastructure Layer** (`db/`, `middleware/`, `config/`) - Provides technical capabilities 

## Architectural Changes

### MongoDB Integration

The application now uses MongoDB as its primary data store with the following components:

1. **Connection Manager**: Handles connections to MongoDB with error handling, retries, and connection pooling.
2. **Repository Pattern**: Abstracts database operations behind repository interfaces.
3. **Data Models**: Domain models with proper BSON/JSON mappings.
4. **DTO Pattern**: Data Transfer Objects for API requests/responses.
5. **Service Layer**: Business logic abstraction over repositories.
6. **Optimized Indexes**: Performance-tuned indexes for common queries.

### Database Layer

The repository pattern is now fully implemented:

```go
// Repository interfaces
type ChatRepository interface {
    Create(ctx context.Context, chat *models.Chat) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*models.Chat, error)
    FindAll(ctx context.Context, limit, offset int) ([]*models.Chat, error)
    Update(ctx context.Context, chat *models.Chat) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    IncrementMessageCount(ctx context.Context, id primitive.ObjectID) error
    CountAll(ctx context.Context) (int64, error)
}

type MessageRepository interface {
    Create(ctx context.Context, message *models.Message) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*models.Message, error)
    FindByChatID(ctx context.Context, chatID primitive.ObjectID, limit, offset int) ([]*models.Message, error)
    CountByChatID(ctx context.Context, chatID primitive.ObjectID) (int64, error)
    Delete(ctx context.Context, id primitive.ObjectID) error
    DeleteByChatID(ctx context.Context, chatID primitive.ObjectID) error
}
```

### Service Layer

The service layer implements business logic and handles the conversion between string IDs and ObjectIDs:

```go
// Chat service example
type ChatService interface {
    CreateChat(ctx context.Context, title string) (*models.Chat, error)
    GetChatByID(ctx context.Context, id string) (*models.Chat, error)
    ListChats(ctx context.Context, page, pageSize int) ([]*models.Chat, int64, error)
    UpdateChat(ctx context.Context, id string, title string) (*models.Chat, error)
    DeleteChat(ctx context.Context, id string) error
}
```

### Key Components

#### HTTP Server (Gin)

The application uses Gin as the HTTP server framework for handling routes and middleware.

```go
// API route groups organization
apiV1 := router.Group("/api/v1")
{
    // Chat routes
    chatRoutes := apiV1.Group("/chats")
    {
        chatRoutes.GET("", handler.ListChats)
        chatRoutes.POST("", handler.CreateChat)
        chatRoutes.GET("/:id", handler.GetChat)
        chatRoutes.PUT("/:id", handler.UpdateChat)
        chatRoutes.DELETE("/:id", handler.DeleteChat)

        // Message routes (nested under chat)
        chatRoutes.GET("/:id/messages", handler.GetMessages)
        chatRoutes.POST("/:id/messages", handler.CreateMessage)
    }

    // Individual message routes
    messages := apiV1.Group("/messages")
    {
        messages.GET("/:id", handler.GetMessage)
        messages.DELETE("/:id", handler.DeleteMessage)
    }
}
```

#### DTO Pattern

The Data Transfer Object pattern separates API contracts from domain models:

```go
// Example DTO for chat responses
type ChatResponse struct {
    ID            string `json:"id"`
    Title         string `json:"title"`
    CreatedAt     string `json:"created_at"`
    UpdatedAt     string `json:"updated_at"`
    LastMessageAt string `json:"last_message_at,omitempty"`
    MessageCount  int    `json:"message_count"`
}
```

#### Middleware Stack

The application uses a layered middleware approach:

```go
// Middleware stack
router.Use(middleware.CORSMiddleware(cfg.Server.AllowedOrigins))
router.Use(middleware.LoggerMiddleware())
router.Use(middleware.ErrorHandlerMiddleware())
```

## Data Flow

1. Client sends a request to a handler endpoint
2. Handler parses the request and converts it to domain parameters
3. Service layer performs business logic and calls repositories
4. Repository interacts with MongoDB
5. Service processes the results
6. Handler converts domain objects to DTOs and returns response

## Configuration

The application uses environment variables for configuration, which can be set via a `.env` file, environment variables, or command-line flags:

```
# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017/sse-chat
MONGODB_DATABASE=sse-chat
MONGODB_TIMEOUT=10s
MONGODB_MAX_POOL_SIZE=100
MONGODB_CONNECT_RETRY_COUNT=5
MONGODB_CONNECT_RETRY_DELAY=3s
```

## Deployment Options

The application can be deployed in several ways:

1. **Docker** - Using the provided Dockerfile and docker-compose.yml
2. **Kubernetes** - For scalable deployments (configuration in k8s/ directory)
3. **Traditional deployment** - Running the binary directly on a server