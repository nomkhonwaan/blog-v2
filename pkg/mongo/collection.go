package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a wrapped interface to the `mongo.Collection` for testing purpose
type Collection interface {
	// Perform finding the documents matching a model
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
}
