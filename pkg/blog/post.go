package blog

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

// Post is a piece of content in the blog platform
type Post struct {
	// Identifier of the post
	ID primitive.ObjectID `bson:"_id" json:"id"`

	// Title of the post
	Title string `bson:"title" json:"title"`

	// Valid URL string composes with title and ID
	Slug string `bson:"slug" json:"slug"`

	// Status of the post which could be...
	// - PUBLISHED
	// - DRAFT
	Status `bson:"status" json:"status"`

	// Original content of the post in markdown syntax
	Markdown string `bson:"markdown" json:"markdown"`

	// Content of the post in HTML format which will be translated from markdown
	HTML string `bson:"html" json:"html"`

	// Date-time that the post was published
	PublishedAt time.Time `bson:"publishedAt" json:"publishedAt"`

	// Identifier of the author
	AuthorID string

	// List of categories that the post belonging to
	Categories []Category

	// List of tags that the post belonging to
	Tags []Tag

	// Date-time that the post was created
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`

	// Date-time that the post was updated
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// MarshalJSON is a custom JSON marshaling function of post entity
func (p *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    p.ID.Hex(),
		Alias: (*Alias)(p),
	})
}

// PostRepository is a repository interface of post
// which defines all post entity related functions
type PostRepository interface {
	// Returns list of posts that filtered by status
	FindAll(FindAllPostsQueryBuilder) ([]Post, error)

	// Returns a single post by its ID
	FindByID(id string) (Post, error)
}

// NewPostRepository returns post repository which connects to MongoDB
func NewPostRepository(col *mongo.Collection) PostRepository {
	return postRepository{col}
}

type postRepository struct {
	col *mongo.Collection
}

func (repo postRepository) FindAll(qb FindAllPostsQueryBuilder) ([]Post, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	filter := bson.M{}
	opts := &options.FindOptions{}

	if qb.Status() != "" {
		filter["status"] = qb.Status()

		if qb.Status() == Published {
			opts.Sort = map[string]interface{}{
				"publishedAt": -1,
			}
		}
	}

	opts.SetSkip(qb.Offset()).SetLimit(qb.Limit())

	cur, err := repo.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	posts := make([]Post, 0)
	for cur.Next(ctx) {
		var p Post

		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (repo postRepository) FindByID(id string) (Post, error) {
	return Post{}, nil
}

type FindAllPostsQueryBuilder interface {
	QueryBuilder

	// Returns status of the post to-be filtered with
	Status() Status

	// Allows to filter posts by status
	WithStatus(Status) FindAllPostsQueryBuilder
}

// NewFindAllPostsQueryBuilder returns new find all posts query builder object with default values
func NewFindAllPostsQueryBuilder() FindAllPostsQueryBuilder {
	return &findAllPostsQueryBuilder{
		QueryBuilder: NewQueryBuilder(),
	}
}

type findAllPostsQueryBuilder struct {
	QueryBuilder

	// Status of the post
	status Status
}

func (qb *findAllPostsQueryBuilder) Status() Status {
	return qb.status
}

func (qb *findAllPostsQueryBuilder) WithStatus(status Status) FindAllPostsQueryBuilder {
	qb.status = status
	return qb
}

// MakePostsHandler creates a new HTTP handler for the "posts" resource
func MakePostsHandler(service Service) http.Handler {
	findAllPublishedPostsHandler := kithttp.NewServer(
		makeFindAllPublishedPostsEndpoint(service),
		decodeFindAllPublishedPostsRequest,
		encodeResponse,
	)

	r := mux.NewRouter().PathPrefix("/v1/posts").Subrouter()

	r.Handle("", findAllPublishedPostsHandler).Methods("GET")

	return r
}

type findAllPublishedPosts struct {
	limit  int64
	offset int64
}

func decodeFindAllPublishedPostsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	queryParams := r.URL.Query()
	offset, limit := queryParams.Get("offset"), queryParams.Get("limit")

	o, _ := strconv.Atoi(offset)
	l, _ := strconv.Atoi(limit)

	if l == 0 {
		l = 5
	}

	return findAllPublishedPosts{
		offset: int64(o),
		limit:  int64(l),
	}, nil
}

func makeFindAllPublishedPostsEndpoint(service Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(findAllPublishedPosts)

		NewFindAllPostsQueryBuilder()

		qb := NewFindAllPostsQueryBuilder()

		qb.WithStatus(Published)
		qb.WithOffset(req.offset).WithLimit(req.limit)

		return service.Post().FindAll(qb)
	}
}

type findPostByIDRequest struct {
	id string
}

func decodeFindPostByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return findPostByIDRequest{""}, nil
}

func makeFindPostByIDEndpoint(service service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		//r := request.(findPostByIDRequest)

		return nil, nil
	}
}
