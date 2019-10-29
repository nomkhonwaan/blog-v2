package blog

import (
	"context"
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// CategoryRepository is a repository interface of category which defines all category entity related functions
type CategoryRepository interface {
	// FindAll returns list of categories
	FindAll(ctx context.Context) ([]Category, error)

	// FindAllByIDs returns list of categories from list of IDs
	FindAllByIDs(ctx context.Context, ids interface{}) ([]Category, error)

	// FindByID returns a single category from its ID
	FindByID(ctx context.Context, id interface{}) (Category, error)
}

// NewCategoryRepository returns category repository
func NewCategoryRepository(col mongo.Collection) MongoCategoryRepository {
	return MongoCategoryRepository{col}
}

// MongoCategoryRepository is a MongoDB specified repository for category
type MongoCategoryRepository struct {
	col mongo.Collection
}

func (repo MongoCategoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	cur, err := repo.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var categories []Category
	err = cur.Decode(&categories)

	return categories, err
}

func (repo MongoCategoryRepository) FindAllByIDs(ctx context.Context, ids interface{}) ([]Category, error) {
	cur, err := repo.col.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids.([]primitive.ObjectID),
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var categories []Category
	err = cur.Decode(&categories)

	return categories, err
}

func (repo MongoCategoryRepository) FindByID(ctx context.Context, id interface{}) (Category, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var cat Category
	err := r.Decode(&cat)

	return cat, err
}
