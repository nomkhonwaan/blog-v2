//go:generate mockgen -destination=./mock/collection_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Collection

package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a wrapped interface to the original mongo.Collection for testing benefit
type Collection interface {
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type collection struct{ *mongo.Collection }

// NewCollection returns new collection which embeds `mongo.Collection` inside
func NewCollection(col *mongo.Collection) Collection {
	return collection{col}
}

// Find executes a find command and returns a Cursor over the matching documents in the collection
func (col collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	cur, err := col.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return cursor{Context: ctx, Cursor: cur}, nil
}

// FindOne executes a find command and returns a SingleResult for one document in the collection
func (col collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	return col.Collection.FindOne(ctx, filter, opts...)
}
