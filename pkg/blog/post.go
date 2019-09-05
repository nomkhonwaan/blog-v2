package blog

import (
	"context"
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Post is a piece of content in the blog platform
type Post struct {
	// Identifier of the post
	ID primitive.ObjectID `bson:"_id" json:"id" graphql:"-"`

	// Title of the post
	Title string `bson:"title" json:"title" graphql:"title"`

	// Valid URL string composes with title and ID
	Slug string `bson:"slug" json:"slug" graphql:"slug"`

	// Status of the post which could be...
	// - PUBLISHED
	// - DRAFT
	Status Status `bson:"status" json:"status" graphql:"status"`

	// Original content of the post in markdown syntax
	Markdown string `bson:"markdown" json:"markdown" graphql:"markdown"`

	// Content of the post in HTML format which will be translated from markdown
	HTML string `bson:"html" json:"html" graphql:"html"`

	// Date-time that the post was published
	PublishedAt time.Time `bson:"publishedAt" json:"publishedAt" graphql:"publishedAt"`

	// Identifier of the author
	AuthorID string `graphql:"-"`

	// List of categories (in reference to the category collection) that the post belonging to
	DBRefCategories []mongo.DBRef `bson:"categories" json:"-" graphql:"-"`

	// List of tags that the post belonging to
	Tags []Tag `graphql:"-"`

	// Date-time that the post was created
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`

	// Date-time that the post was updated
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt" graphql:"updatedAt"`
}

// MarshalJSON is a custom JSON marshaling function of post entity
func (p Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    p.ID.Hex(),
		Alias: (*Alias)(&p),
	})
}

func (Post) Categories(repository CategoryRepository) interface{} {
	return func(ctx context.Context, p Post) ([]Category, error) {
		ids := make([]primitive.ObjectID, len(p.DBRefCategories))
		for i, dbRef := range p.DBRefCategories {
			ids[i] = dbRef.ID
		}
		return repository.FindAllByIDs(ctx, ids)
	}
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
func NewPostRepository(col mongo.Collection) MongoPostRepository {
	return MongoPostRepository{col}
}

type MongoPostRepository struct {
	col mongo.Collection
}

func (repo MongoPostRepository) FindAll(ctx context.Context, q PostQuery) ([]Post, error) {
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

	return repo.scanAll(ctx, cur)
}

func (repo MongoPostRepository) FindByID(ctx context.Context, id string) (Post, error) {
	return Post{}, nil
}

func (repo MongoPostRepository) scanAll(ctx context.Context, cur *mgo.Cursor) ([]Post, error) {
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
