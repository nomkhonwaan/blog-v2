package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect is a wrapped function to the original mongo.Connect for avoiding package name conflict
func Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts...)
}
