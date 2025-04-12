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
	"github.com/yasin-yalcin-dev/go-sse-ai-chat/internal/config"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// createClientOptions creates MongoDB client options from configuration
func createClientOptions(cfg *config.MongoDBConfig) *options.ClientOptions {
	return options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetConnectTimeout(cfg.Timeout)
}
