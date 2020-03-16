package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a wrapped interface to the original mongo.Database for testing benefit
type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}
