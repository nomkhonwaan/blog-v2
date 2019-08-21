package blog

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
}

// NewCategoryRepository returns category repository which connects to MongoDB
func NewCategoryRepository(col *mongo.Collection) CategoryRepository {
	return categoryRepository{col}
}

type categoryRepository struct {
	col *mongo.Collection
}

func (repo categoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	cur, err := repo.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

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

// MakeCategoriesHandler creates a new HTTP handler for the "categories" resource
func MakeCategoriesHandler(service Service) http.Handler {
	findAllCategoriesHandler := kithttp.NewServer(
		makeFindAllCategoriesEndpoint(service),
		decodeFindAllCategoriesRequest,
		encodeResponse,
	)

	r := mux.NewRouter().PathPrefix("/v1/categories").Subrouter()

	r.Handle("", findAllCategoriesHandler).Methods("GET")

	return r
}

type findAllCategoriesRequest struct{}

func decodeFindAllCategoriesRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return findAllCategoriesRequest{}, nil
}

func makeFindAllCategoriesEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(findAllCategoriesRequest)

		return service.Category().FindAll(ctx)
	}
}
