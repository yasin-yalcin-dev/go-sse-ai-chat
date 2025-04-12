# Setup and Installation Guide

This document provides comprehensive instructions for setting up the Go SSE AI Chat project for development and production use.

## Prerequisites

Before starting, ensure you have the following installed:

- Go 1.20 or higher
- MongoDB 5.0 or higher
- Docker and Docker Compose (optional)
- An OpenAI or Anthropic API key

## Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yasin-yalcin-dev/go-sse-ai-chat.git
cd go-sse-ai-chat
```

### 2. Set Up Environment Variables

Create a `.env` file in the project root:

```bash
# Server Configuration
SERVER_PORT=8080
ENVIRONMENT=development

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017/sse-chat
MONGODB_DATABASE=sse-chat

# AI Provider Configuration
AI_PROVIDER=openai  # Options: openai, anthropic
OPENAI_API_KEY=your_openai_api_key
ANTHROPIC_API_KEY=your_anthropic_api_key

# SSE Configuration
SSE_MAX_CLIENTS=1000
SSE_KEEPALIVE_INTERVAL=15s
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run MongoDB

If you don't have MongoDB running locally, you can use Docker:

```bash
docker run --name mongodb -d -p 27017:27017 mongo:latest
```

### 5. Build and Run the Application

```bash
go build -o sse-chat ./cmd/server
./sse-chat
```

Alternatively, use Go run:

```bash
go run ./cmd/server/main.go
```

### 6. Access the Application

Open your browser and navigate to:
- Frontend: `http://localhost:8080`
- API Health Check: `http://localhost:8080/health`

## Docker Development Setup

### 1. Using Docker Compose

```bash
# Start the application and MongoDB
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 2. Build the Docker Image Manually

```bash
docker build -t go-sse-ai-chat:latest .
```

## Production Deployment

### 1. Environment Configuration

For production, ensure you set the following environment variables:

```bash
ENVIRONMENT=production
SERVER_PORT=8080
MONGODB_URI=your_production_mongodb_uri
OPENAI_API_KEY=your_production_api_key
```

### 2. Docker Production Deployment

```bash
docker run -d --name go-sse-ai-chat \
  -p 8080:8080 \
  -e ENVIRONMENT=production \
  -e MONGODB_URI=your_production_mongodb_uri \
  -e OPENAI_API_KEY=your_production_api_key \
  go-sse-ai-chat:latest
```

### 3. Kubernetes Deployment (Optional)

Apply the Kubernetes manifests:

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
```

## Configuration Options

### Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port to run the server on | 8080 |
| ENVIRONMENT | Application environment | development |
| LOG_LEVEL | Logging level | info |

### MongoDB Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| MONGODB_URI | MongoDB connection URI | mongodb://localhost:27017/sse-chat |
| MONGODB_DATABASE | MongoDB database name | sse-chat |
| MONGODB_TIMEOUT | Connection timeout in seconds | 10 |

### AI Provider Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| AI_PROVIDER | AI provider to use (openai or anthropic) | openai |
| OPENAI_API_KEY | OpenAI API key | - |
| OPENAI_MODEL | OpenAI model to use | gpt-4o |
| ANTHROPIC_API_KEY | Anthropic API key | - |
| ANTHROPIC_MODEL | Anthropic model to use | claude-3-opus-20240229 |

### SSE Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| SSE_MAX_CLIENTS | Maximum number of concurrent SSE clients | 1000 |
| SSE_KEEPALIVE_INTERVAL | Interval for sending keepalive messages | 15s |
| SSE_RECONNECT_TIMEOUT | Time window for client reconnection | 1m |

## Troubleshooting

### Common Issues

1. **Connection refused to MongoDB**
   - Check if MongoDB is running: `docker ps | grep mongo`
   - Verify the connection string in your .env file

2. **API Key Authentication Errors**
   - Ensure your OpenAI/Anthropic API key is valid
   - Check if you've set the correct provider in AI_PROVIDER

3. **Port Already in Use**
   - Change the SERVER_PORT in your .env file
   - Check for other processes using port 8080: `lsof -i :8080`

### Getting Help

If you encounter issues not covered in this guide, please:
1. Check the project issues on GitHub
2. Create a new issue with detailed information about your problem
