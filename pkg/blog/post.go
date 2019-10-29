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

// PostRepository is a repository interface of post which defines all post entity related functions
type PostRepository interface {
	// Create inserts a new empty post which belongs to the author with "Draft" status
	Create(ctx context.Context, authorID string) (Post, error)

	// FindAll returns list of posts filtered by post query
	FindAll(ctx context.Context, q PostQuery) ([]Post, error)

	// FindByID returns a single post from its ID
	FindByID(ctx context.Context, id interface{}) (Post, error)

	// Save does updating a single post
	Save(ctx context.Context, id interface{}, q PostQuery) (Post, error)
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
	if cat := q.Category(); !cat.ID.IsZero() {
		filter["categories.$id"] = cat.ID
	}
	if tag := q.Tag(); !tag.ID.IsZero() {
		filter["tags.$id"] = tag.ID
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

func (repo MongoPostRepository) Save(ctx context.Context, id interface{}, q PostQuery) (Post, error) {
	update := bson.M{"$set": bson.M{}}

	if title := q.Title(); title != "" {
		update["$set"].(bson.M)["title"] = title
	}
	if markdown := q.Markdown(); markdown != "" {
		update["$set"].(bson.M)["markdown"] = markdown
	}
	if html := q.HTML(); html != "" {
		update["$set"].(bson.M)["html"] = html
	}
	if categories := q.Categories(); categories != nil {
		update["$set"].(bson.M)["categories"] = make(primitive.A, 0)
		for _, cat := range categories {
			update["$set"].(bson.M)["categories"] = append(
				update["$set"].(bson.M)["categories"].(primitive.A),
				mongo.DBRef{
					Ref: "categories",
					ID:  cat.ID,
				},
			)
		}
	}
	if tags := q.Tags(); tags != nil {
		update["$set"].(bson.M)["tags"] = make(primitive.A, 0)
		for _, tag := range tags {
			update["$set"].(bson.M)["tags"] = append(
				update["$set"].(bson.M)["tags"].(primitive.A),
				mongo.DBRef{
					Ref: "tags",
					ID:  tag.ID,
				},
			)
		}
	}

	_, err := repo.col.UpdateOne(ctx, bson.M{"_id": id.(primitive.ObjectID)}, update)
	if err != nil {
		return Post{}, err
	}

	return repo.FindByID(ctx, id.(primitive.ObjectID))
}

// PostQueryBuilder is a builder for building query object that repository can use to find all posts
type PostQueryBuilder interface {
	// WithTitle allows to set title to the post query object
	WithTitle(title string) PostQueryBuilder

	//WithMarkdown allows to set markdown to the post query object
	WithMarkdown(markdown string) PostQueryBuilder

	// WithHTML allows to set HTML to the post query object
	WithHTML(html string) PostQueryBuilder

	// WithStatus allows to set status to the post query object
	WithStatus(status Status) PostQueryBuilder

	// WithCategory allows to set category to the post query object
	WithCategory(category Category) PostQueryBuilder

	// WithCategories allows to set categories to the post query object
	WithCategories(categories []Category) PostQueryBuilder

	// WithTag allows to set tag to the post query object
	WithTag(tag Tag) PostQueryBuilder

	// WithTags allows to set tags to the post query object
	WithTags(tags []Tag) PostQueryBuilder

	// WithOffset allows to set returned result offset
	WithOffset(offset int64) PostQueryBuilder

	// WithLimit allows to set maximum returned result
	WithLimit(limit int64) PostQueryBuilder

	// Build returns a post query object
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

func (qb *MongoPostQueryBuilder) WithTitle(title string) PostQueryBuilder {
	qb.MongoPostQuery.title = title
	return qb
}

func (qb *MongoPostQueryBuilder) WithMarkdown(markdown string) PostQueryBuilder {
	qb.MongoPostQuery.markdown = markdown
	return qb
}

func (qb *MongoPostQueryBuilder) WithHTML(html string) PostQueryBuilder {
	qb.MongoPostQuery.html = html
	return qb
}

func (qb *MongoPostQueryBuilder) WithStatus(status Status) PostQueryBuilder {
	qb.MongoPostQuery.status = status
	return qb
}

func (qb *MongoPostQueryBuilder) WithCategory(category Category) PostQueryBuilder {
	qb.MongoPostQuery.category = category
	return qb
}

func (qb *MongoPostQueryBuilder) WithCategories(categories []Category) PostQueryBuilder {
	qb.MongoPostQuery.categories = categories
	return qb
}

func (qb *MongoPostQueryBuilder) WithTag(tag Tag) PostQueryBuilder {
	qb.MongoPostQuery.tag = tag
	return qb
}

func (qb *MongoPostQueryBuilder) WithTags(tags []Tag) PostQueryBuilder {
	qb.MongoPostQuery.tags = tags
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
	// Title returns a title field
	Title() string

	// Markdown returns markdown field
	Markdown() string

	// HTML returns an HTML field
	HTML() string

	// Status returns a status field
	Status() Status

	// Category returns a single category field
	Category() Category

	// Categories returns a list of categories field
	Categories() []Category

	// Tag returns a single tag field
	Tag() Tag

	// Tags returns a list of tags field
	Tags() []Tag

	// Offset returns an offset of the returned result
	Offset() int64

	// Limit returns  maximum number of the returned result
	Limit() int64
}

type MongoPostQuery struct {
	title      string
	markdown   string
	html       string
	status     Status
	category   Category
	categories []Category
	tag        Tag
	tags       []Tag

	offset int64
	limit  int64
}

func (q *MongoPostQuery) Title() string {
	return q.title
}

func (q *MongoPostQuery) Markdown() string {
	return q.markdown
}

func (q *MongoPostQuery) HTML() string {
	return q.html
}

func (q *MongoPostQuery) Status() Status {
	return q.status
}

func (q *MongoPostQuery) Category() Category {
	return q.category
}

func (q *MongoPostQuery) Categories() []Category {
	return q.categories
}

func (q *MongoPostQuery) Tag() Tag {
	return q.tag
}

func (q *MongoPostQuery) Tags() []Tag {
	return q.tags
}

func (q *MongoPostQuery) Offset() int64 {
	return q.offset
}

func (q *MongoPostQuery) Limit() int64 {
	return q.limit
}
