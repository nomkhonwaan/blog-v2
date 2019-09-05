package blog

import (
	"context"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
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

// CategoryRepository is a repository interface of category
// which defines all category entity related functions
type CategoryRepository interface {
	// Returns list of categories
	FindAll(ctx context.Context) ([]Category, error)

	// Returns list of categories from list of IDs
	FindAllByIDs(ctx context.Context, ids []primitive.ObjectID) ([]Category, error)
}

// NewCategoryRepository returns category repository which connects to MongoDB
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

	return repo.scanAll(ctx, cur)
}

func (repo MongoCategoryRepository) FindAllByIDs(ctx context.Context, ids []primitive.ObjectID) ([]Category, error) {
	cur, err := repo.col.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return repo.scanAll(ctx, cur)
}

func (repo MongoCategoryRepository) scanAll(ctx context.Context, cur *mgo.Cursor) ([]Category, error) {
	categories := make([]Category, 0)

	for cur.Next(ctx) {
		var c Category

		err := cur.Decode(&c)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}
