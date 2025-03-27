/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config is the main configuration structure
type Config struct {
	Server     ServerConfig
	MongoDB    MongoDBConfig
	SSE        SSEConfig
	LogLevel   string
	AIProvider AIProviderConfig
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	ShutdownTimeout  time.Duration
	RequestBodyLimit int64
	TrustedProxies   string
	AllowedOrigins   []string
	DefaultPageSize  int
	MaxPageSize      int
}

// MongoDBConfig contains MongoDB configuration
type MongoDBConfig struct {
	URI               string
	Database          string
	Timeout           time.Duration
	MaxPoolSize       uint64
	ConnectRetryCount int
	ConnectRetryDelay time.Duration
}

// SSEConfig contains Server-Sent Events configuration
type SSEConfig struct {
	MaxClients        int
	KeepaliveInterval time.Duration
	BufferSize        int
	WriteTimeout      time.Duration
}

// AIProviderConfig contains AI provider configuration
type AIProviderConfig struct {
	Provider       string // "openai" or "anthropic"
	OpenAIKey      string
	OpenAIModel    string
	AnthropicKey   string
	AnthropicModel string
	Timeout        time.Duration
	MaxTokens      int
}

// Load Loads the .env file and environment variables
func Load() (*Config, error) {
	// Load .env file, otherwise use environment variables
	godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:             getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:      getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:     getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout:  getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 15*time.Second),
			RequestBodyLimit: int64(getEnvInt("SERVER_REQUEST_BODY_LIMIT", 1024)) * 1024, // KB -> Bytes
			TrustedProxies:   getEnv("SERVER_TRUSTED_PROXIES", "127.0.0.1"),
			AllowedOrigins:   getEnvSlice("SERVER_ALLOWED_ORIGINS", []string{"*"}),
			DefaultPageSize:  getEnvInt("SERVER_DEFAULT_PAGE_SIZE", 20),
			MaxPageSize:      getEnvInt("SERVER_MAX_PAGE_SIZE", 100),
		},
		MongoDB: MongoDBConfig{
			URI:               getEnv("MONGODB_URI", "mongodb://localhost:27017/sse-chat"),
			Database:          getEnv("MONGODB_DATABASE", "sse-chat"),
			Timeout:           getEnvDuration("MONGODB_TIMEOUT", 10*time.Second),
			MaxPoolSize:       uint64(getEnvInt("MONGODB_MAX_POOL_SIZE", 100)),
			ConnectRetryCount: getEnvInt("MONGODB_CONNECT_RETRY_COUNT", 5),
			ConnectRetryDelay: getEnvDuration("MONGODB_CONNECT_RETRY_DELAY", 3*time.Second),
		},
		SSE: SSEConfig{
			MaxClients:        getEnvInt("SSE_MAX_CLIENTS", 1000),
			KeepaliveInterval: getEnvDuration("SSE_KEEPALIVE_INTERVAL", 15*time.Second),
			BufferSize:        getEnvInt("SSE_BUFFER_SIZE", 256),
			WriteTimeout:      getEnvDuration("SSE_WRITE_TIMEOUT", 5*time.Second),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		AIProvider: AIProviderConfig{
			Provider:       getEnv("AI_PROVIDER", "openai"),
			OpenAIKey:      getEnv("OPENAI_API_KEY", ""),
			OpenAIModel:    getEnv("OPENAI_MODEL", "gpt-4o"),
			AnthropicKey:   getEnv("ANTHROPIC_API_KEY", ""),
			AnthropicModel: getEnv("ANTHROPIC_MODEL", "claude-3-opus-20240229"),
			Timeout:        getEnvDuration("AI_TIMEOUT", 60*time.Second),
			MaxTokens:      getEnvInt("AI_MAX_TOKENS", 4096),
		},
	}

	// verify configuration
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that the configuration is valid
func validate(cfg *Config) error {
	//Server control
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("SERVER_PORT invalid: %d", cfg.Server.Port)
	}

	// MongoDB control
	if cfg.MongoDB.URI == "" {
		return fmt.Errorf("MONGODB_URI is required")
	}

	// AI Provider control
	provider := cfg.AIProvider.Provider
	if provider != "openai" && provider != "anthropic" {
		return fmt.Errorf("AI_PROVIDER value must be 'openai' or 'anthropic', received: %s", provider)
	}

	if provider == "openai" && cfg.AIProvider.OpenAIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required when AI_PROVIDER is set to 'openai'")
	}

	if provider == "anthropic" && cfg.AIProvider.AnthropicKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY is required when AI_PROVIDER is set to 'anthropic'")
	}

	return nil
}

// Auxiliary functions

// getEnv takes an environment variable key and a default value, and returns the value of the environment variable or the default value if it is not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt takes an environment variable key and a default value, and returns the value of the environment variable or the default value if it is not set
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration takes an environment variable key and a default value, and returns the value of the environment variable or the default value if it is not set
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvBool takes an environment variable key and a default value, and returns the value of the environment variable or the default value if it is not set
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvSlice takes an environment variable key and a default value, and returns the value of the environment variable or the default value if it is not set
func getEnvSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		// Split comma separated values
		// TODO: A better parsing algorithm can be added (like CSV)
		return splitCSV(value)
	}
	return defaultValue
}

// splitCSV splits a string by commas
func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}

	// TODO: Implement proper CSV splitting
	return []string{s}
}
