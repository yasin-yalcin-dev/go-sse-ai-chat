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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client returns the MongoDB client
func (c *DBConnection) Client() *mongo.Client {
	return c.client
}

// Database returns the MongoDB database
func (c *DBConnection) Database() *mongo.Database {
	return c.database
}

// Chats returns the chats collection
func (c *DBConnection) Chats() *mongo.Collection {
	return c.database.Collection(c.cfg.CollectionChats)
}

// Messages returns the messages collection
func (c *DBConnection) Messages() *mongo.Collection {
	return c.database.Collection(c.cfg.CollectionMessages)
}

// Collection returns a MongoDB collection
func (c *DBConnection) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// CollectionWithOptions returns a MongoDB collection with specific options
func (c *DBConnection) CollectionWithOptions(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return c.database.Collection(name, opts...)
}
