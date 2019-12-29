//go:generate mockgen -destination=./mock/tag_mock.go github.com/nomkhonwaan/myblog/pkg/blog TagRepository

package blog

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag is a label attached to the post for the purpose of identification
type Tag struct {
	// Identifier of the tag
	ID primitive.ObjectID `bson:"_id" json:"id" graphql:"-"`
	// Name of the tag
	Name string `bson:"name" json:"name" graphql:"name"`
	// Valid URL string composes with name and ID
	Slug string `bson:"slug" json:"slug" graphql:"slug"`
}

// MarshalJSON is a custom JSON marshaling function of tag entity
func (tag Tag) MarshalJSON() ([]byte, error) {
	type Alias Tag
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    tag.ID.Hex(),
		Alias: (*Alias)(&tag),
	})
}

// TagRepository is a repository interface of category which defines all category entity related functions
type TagRepository interface {
	// FindAll returns list of tags
	FindAll(ctx context.Context) ([]Tag, error)

	// FindAllByIDs returns list of tags from list of IDs
	FindAllByIDs(ctx context.Context, ids interface{}) ([]Tag, error)

	// FindByID returns a single tag from its ID
	FindByID(ctx context.Context, id interface{}) (Tag, error)
}

// NewTagRepository returns tag repository
func NewTagRepository(col mongo.Collection) MongoTagRepository {
	return MongoTagRepository{col}
}

// MongoTagRepository is a MongoDB specified repository for tag
type MongoTagRepository struct {
	col mongo.Collection
}

func (repo MongoTagRepository) FindAll(ctx context.Context) ([]Tag, error) {
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cur, err := repo.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tags []Tag
	err = cur.Decode(&tags)

	return tags, err
}

func (repo MongoTagRepository) FindAllByIDs(ctx context.Context, ids interface{}) ([]Tag, error) {
	filter := bson.M{"_id": bson.M{"$in": ids.([]primitive.ObjectID)}}
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cur, err := repo.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tags []Tag
	err = cur.Decode(&tags)

	return tags, err
}

func (repo MongoTagRepository) FindByID(ctx context.Context, id interface{}) (Tag, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var tag Tag
	err := r.Decode(&tag)

	return tag, err
}
