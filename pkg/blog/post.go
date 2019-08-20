package blog

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Returns list of posts
	FindAll(ctx context.Context, q PostQuery) ([]Post, error)

	// Returns a single post by its ID
	FindByID(ctx context.Context, id string) (Post, error)
}

// NewPostRepository returns post repository which connects to MongoDB
func NewPostRepository(col *mongo.Collection) PostRepository {
	return postRepository{col}
}

type postRepository struct {
	col *mongo.Collection
}

func (repo postRepository) FindAll(ctx context.Context, q PostQuery) ([]Post, error) {
	filter := bson.M{}
	opts := &options.FindOptions{}

	if q.Status() != "" {
		filter["status"] = q.Status()

		if q.Status() == Published {
			opts.Sort = map[string]interface{}{
				"publishedAt": -1,
			}
		}
	}

	opts.SetSkip(q.Offset()).SetLimit(q.Limit())

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

func (repo postRepository) FindByID(ctx context.Context, id string) (Post, error) {
	return Post{}, nil
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(findAllPublishedPosts)
		q := NewPostQueryBuilder().
			WithStatus(Published).
			WithOffset(req.offset).
			WithLimit(req.limit).
			Build()

		return service.Post().FindAll(ctx, q)
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

// PostQueryBuilder is a builder for building query object that repository can use to find all posts
type PostQueryBuilder interface {
	// Allow to filter post by status
	WithStatus(status Status) PostQueryBuilder

	// Allow to set returned result offset
	WithOffset(offset int64) PostQueryBuilder

	// Allow to set maximum returned result
	WithLimit(limit int64) PostQueryBuilder

	// Return a post query
	Build() PostQuery
}

// NewPostQueryBuilder returns a query builder for building post query object
func NewPostQueryBuilder() PostQueryBuilder {
	return &postQueryBuilder{
		postQuery: &postQuery{
			offset: 0,
			limit:  5,
		},
	}
}

type postQueryBuilder struct {
	*postQuery
}

func (qb *postQueryBuilder) WithStatus(status Status) PostQueryBuilder {
	qb.postQuery.status = status
	return qb
}

func (qb *postQueryBuilder) WithOffset(offset int64) PostQueryBuilder {
	qb.postQuery.offset = offset
	return qb
}

func (qb *postQueryBuilder) WithLimit(limit int64) PostQueryBuilder {
	qb.postQuery.limit = limit
	return qb
}

func (qb *postQueryBuilder) Build() PostQuery {
	return qb.postQuery
}

// PostQuery is a query object which will be used for filtering list of posts
type PostQuery interface {
	// Return status to be filtered with
	Status() Status

	// Return offset of the returned result
	Offset() int64

	// Return maximum number of the returned result
	Limit() int64
}

type postQuery struct {
	status Status
	offset int64
	limit  int64
}

func (q *postQuery) Status() Status {
	return q.status
}

func (q *postQuery) Offset() int64 {
	return q.offset
}

func (q *postQuery) Limit() int64 {
	return q.limit
}
