package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a wrapped interface to the `mongo.Collection` for testing purpose
type Collection interface {
	// Perform finding the documents matching a model
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error)

	// Perform finding up to one document that matches the model
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult

	// Insert a single document into the collection
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// Update a single document int the collection
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// CustomCollection provides customized collection methods on top of the original `mongo.Collection`
type CustomCollection struct {
	*mongo.Collection
}

// NewCustomCollection returns new CustomCollection which embeds `mongo.Collection` inside
func NewCustomCollection(col *mongo.Collection) CustomCollection {
	return CustomCollection{col}
}

func (col CustomCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	cur, err := col.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return CustomCursor{Context: ctx, Cursor: cur}, nil
}

func (col CustomCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	return col.Collection.FindOne(ctx, filter, opts...)
}
