package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	AuthorID string `bson:"authorId" json:"authorId" graphql:"authorId"`

	// List of categories (in reference to the category collection) that the post belonging to
	Categories []mongo.DBRef `bson:"categories" json:"-" graphql:"-"`

	// List of tags that the post belonging to
	Tags []mongo.DBRef `bson:"tags" json:"-" graphql:"-"`

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

// BelongToCategories returns list of categories that the post has belonging to
func (Post) BelongToCategories(repo CategoryRepository) interface{} {
	return func(ctx context.Context, p Post) ([]Category, error) {
		ids := make([]primitive.ObjectID, len(p.Categories))
		for i, category := range p.Categories {
			ids[i] = category.ID
		}
		return repo.FindAllByIDs(ctx, ids)
	}
}

// BelongToTags returns list of tags that the post has belonging to
func (Post) BelongToTags(repo TagRepository) interface{} {
	return func(ctx context.Context, p Post) ([]Tag, error) {
		ids := make([]primitive.ObjectID, len(p.Tags))
		for i, tag := range p.Tags {
			ids[i] = tag.ID
		}
		return repo.FindAllByIDs(ctx, ids)
	}
}

// PostRepository is a repository interface of post which defines all post entity related functions
type PostRepository interface {
	// Create new empty post which belongs to the author with "Draft" status
	Create(ctx context.Context, authorID string) (Post, error)

	// Return list of posts filtered by post query
	FindAll(ctx context.Context, q PostQuery) ([]Post, error)

	// Return a single post by its ID
	FindByID(ctx context.Context, id interface{}) (Post, error)
}

// NewPostRepository returns post repository
func NewPostRepository(col mongo.Collection) MongoPostRepository {
	return MongoPostRepository{col}
}

// MongoPostRepository is a MongoDB specified repository for post
type MongoPostRepository struct {
	col mongo.Collection
}

func (repo MongoPostRepository) Create(ctx context.Context, authorID string) (Post, error) {
	id := primitive.NewObjectID()
	post := Post{
		ID:        id,
		Slug:      fmt.Sprintf("%s", id.Hex()),
		Status:    Draft,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}

	doc, _ := bson.Marshal(post)
	_, err := repo.col.InsertOne(ctx, doc)
	if err != nil {
		return Post{}, err
	}

	return post, nil
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

	var posts []Post
	err = cur.Decode(&posts)

	return posts, err
}

func (repo MongoPostRepository) FindByID(ctx context.Context, id interface{}) (Post, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var p Post
	err := r.Decode(&p)

	return p, err
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
func NewPostQueryBuilder() *MongoPostQueryBuilder {
	return &MongoPostQueryBuilder{
		MongoPostQuery: &MongoPostQuery{
			offset: 0,
			limit:  5,
		},
	}
}

type MongoPostQueryBuilder struct {
	*MongoPostQuery
}

func (qb *MongoPostQueryBuilder) WithStatus(status Status) PostQueryBuilder {
	qb.MongoPostQuery.status = status
	return qb
}

func (qb *MongoPostQueryBuilder) WithOffset(offset int64) PostQueryBuilder {
	qb.MongoPostQuery.offset = offset
	return qb
}

func (qb *MongoPostQueryBuilder) WithLimit(limit int64) PostQueryBuilder {
	qb.MongoPostQuery.limit = limit
	return qb
}

func (qb *MongoPostQueryBuilder) Build() PostQuery {
	return qb.MongoPostQuery
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

type MongoPostQuery struct {
	status Status
	offset int64
	limit  int64
}

func (q *MongoPostQuery) Status() Status {
	return q.status
}

func (q *MongoPostQuery) Offset() int64 {
	return q.offset
}

func (q *MongoPostQuery) Limit() int64 {
	return q.limit
}
