/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/config"
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DBConnection represents a MongoDB database connection
type DBConnection struct {
	client   *mongo.Client
	database *mongo.Database
	cfg      *config.MongoDBConfig
}

// New creates a new MongoDB connection
func New(cfg *config.MongoDBConfig) (*DBConnection, error) {
	conn := &DBConnection{
		cfg: cfg,
	}

	if err := conn.connect(); err != nil {
		return nil, err
	}

	// Ensure indexes exist
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := conn.EnsureIndexes(ctx); err != nil {
		logger.Warnf("Failed to create indexes: %v", err)
		// Continue even if index creation fails
	}

	return conn, nil
}

// connect establishes a connection to MongoDB
func (c *DBConnection) connect() error {
	// Create MongoDB client options
	clientOptions := createClientOptions(c.cfg)

	// Create connection context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Timeout)
	defer cancel()

	// Connect to MongoDB
	var err error
	var client *mongo.Client

	for attempt := 1; attempt <= c.cfg.ConnectRetryCount; attempt++ {
		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			// Ping the database to verify connection
			pingCtx, pingCancel := context.WithTimeout(context.Background(), c.cfg.Timeout)
			defer pingCancel()

			err = client.Ping(pingCtx, readpref.Primary())
			if err == nil {
				// Connection successful
				c.client = client
				c.database = client.Database(c.cfg.Database)
				logger.Infof("Connected to MongoDB: %s", c.cfg.URI)
				return nil
			}
		}

		logger.Warnf("Failed to connect to MongoDB (attempt %d/%d): %v",
			attempt, c.cfg.ConnectRetryCount, err)

		// Wait before retrying
		if attempt < c.cfg.ConnectRetryCount {
			time.Sleep(c.cfg.ConnectRetryDelay)
		}
	}

	return errors.New("failed to connect to MongoDB after multiple attempts: " + err.Error())
}

// Disconnect closes the MongoDB connection
func (c *DBConnection) Disconnect(ctx context.Context) error {
	if c.client != nil {
		err := c.client.Disconnect(ctx)
		if err != nil {
			return err
		}
		logger.Info("Disconnected from MongoDB")
	}
	return nil
}

// Ping checks the connection to the database
func (c *DBConnection) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}
