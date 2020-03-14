//go:generate mockgen -destination=./mock/category_mock.go github.com/nomkhonwaan/myblog/pkg/blog CategoryRepository

package blog

import (
	"context"
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Category is a group of posts regarded as having particular shared characteristics
type Category struct {
	// Identifier of the category
	ID primitive.ObjectID `bson:"_id" json:"id" graphql:"-"`

	// Name of the category
	Name string `bson:"name" json:"name" graphql:"name"`

	// Valid URL string composes with name and ID
	Slug string `bson:"slug" json:"slug" graphql:"slug"`
}

// MarshalJSON is a custom JSON marshaling function of category entity
func (cat Category) MarshalJSON() ([]byte, error) {
	type Alias Category
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    cat.ID.Hex(),
		Alias: (*Alias)(&cat),
	})
}

// A CategoryRepository interface
type CategoryRepository interface {
	FindAll(ctx context.Context) ([]Category, error)
	FindAllByIDs(ctx context.Context, ids interface{}) ([]Category, error)
	FindByID(ctx context.Context, id interface{}) (Category, error)
}

// NewCategoryRepository returns a MongoCategoryRepository instance
func NewCategoryRepository(col mongo.Collection) MongoCategoryRepository {
	return MongoCategoryRepository{col}
}

// MongoCategoryRepository implements CategoryRepository interface
type MongoCategoryRepository struct{ col mongo.Collection }

// FindAll returns list of categories
func (repo MongoCategoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cur, err := repo.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var categories []Category
	err = cur.Decode(&categories)

	return categories, err
}

// FindAllByIDs returns list of categories from list of IDs
func (repo MongoCategoryRepository) FindAllByIDs(ctx context.Context, ids interface{}) ([]Category, error) {
	filter := bson.M{"_id": bson.M{"$in": ids.([]primitive.ObjectID)}}
	opts := options.Find().SetSort(bson.D{{"name", 1}})
	cur, err := repo.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var categories []Category
	err = cur.Decode(&categories)

	return categories, err
}

// FindByID returns a single category from its ID
func (repo MongoCategoryRepository) FindByID(ctx context.Context, id interface{}) (Category, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var cat Category
	err := r.Decode(&cat)

	return cat, err
}
