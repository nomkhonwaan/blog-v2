//go:generate mockgen -destination=./mock/post_mock.go github.com/nomkhonwaan/myblog/pkg/blog PostRepository

package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
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

	// A featured image to be shown in the social network as a cover image
	FeaturedImage mongo.DBRef `bson:"featuredImage" json:"-" graphql:"-"`

	// List of attachments are belonging to the post
	Attachments []mongo.DBRef `bson:"attachments" json:"-" graphql:"-"`

	// A social network engagement of the post
	Engagement Engagement `bson:"-" json:"engagement" graphql:"engagement"`

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

// A PostRepository interface
type PostRepository interface {
	Create(ctx context.Context, authorID string) (Post, error)
	FindAll(ctx context.Context, q PostQuery) ([]Post, error)
	FindByID(ctx context.Context, id interface{}) (Post, error)
	Save(ctx context.Context, id interface{}, q PostQuery) (Post, error)
}

// NewPostRepository returns a MongoPostRepository instance
func NewPostRepository(db mongo.Database) MongoPostRepository {
	return MongoPostRepository{col: mongo.NewCollection(db.Collection("posts"))}
}

// MongoPostRepository implements PostRepository interface
type MongoPostRepository struct {
	col mongo.Collection
}

// Create inserts a new empty post which belongs to the author with "Draft" status
func (repo MongoPostRepository) Create(ctx context.Context, authorID string) (Post, error) {
	id := primitive.NewObjectID()
	post := Post{
		ID:        id,
		Slug:      fmt.Sprintf("%s", id.Hex()),
		Status:    StatusDraft,
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

// FindAll returns list of posts filtered by post query
func (repo MongoPostRepository) FindAll(ctx context.Context, q PostQuery) ([]Post, error) {
	filter := bson.M{}
	opts := options.Find()

	if status := q.Status(); status != nil {
		filter["status"] = status

		if q.Status().IsPublished() {
			opts.SetSort(bson.D{{"publishedAt", -1}})
		} else if q.Status().IsDraft() {
			opts.SetSort(bson.D{{"createdAt", -1}})
		}
	} else {
		opts.SetSort(bson.D{{"status", 1}, {"createdAt", -1}})
	}
	if authorID := q.AuthorID(); authorID != nil {
		filter["authorId"] = authorID
	}
	if cat := q.Category(); cat != nil {
		filter["categories.$id"] = cat.ID
	}
	if tag := q.Tag(); tag != nil {
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

// FindByID returns a single post from its ID
func (repo MongoPostRepository) FindByID(ctx context.Context, id interface{}) (Post, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var p Post
	err := r.Decode(&p)

	return p, err
}

// Save does updating a single post
func (repo MongoPostRepository) Save(ctx context.Context, id interface{}, q PostQuery) (Post, error) {
	update := bson.M{"$set": bson.M{
		"updatedAt": time.Now(),
	}}

	if title := q.Title(); title != nil {
		update["$set"].(bson.M)["title"] = title
	}
	if slug := q.Slug(); slug != nil {
		update["$set"].(bson.M)["slug"] = slug
	}
	if status := q.Status(); status != nil {
		update["$set"].(bson.M)["status"] = status
	}
	if markdown := q.Markdown(); markdown != nil {
		update["$set"].(bson.M)["markdown"] = markdown
	}
	if html := q.HTML(); html != nil {
		update["$set"].(bson.M)["html"] = html
	}
	if publishedAt := q.PublishedAt(); publishedAt != nil {
		update["$set"].(bson.M)["publishedAt"] = publishedAt
	}
	if categories := q.Categories(); categories != nil {
		update["$set"].(bson.M)["categories"] = make(primitive.A, 0)
		for _, cat := range *categories {
			update["$set"].(bson.M)["categories"] = append(
				update["$set"].(bson.M)["categories"].(primitive.A),
				mongo.DBRef{Ref: "categories", ID: cat.ID},
			)
		}
	}
	if tags := q.Tags(); tags != nil {
		update["$set"].(bson.M)["tags"] = make(primitive.A, 0)
		for _, tag := range *tags {
			update["$set"].(bson.M)["tags"] = append(
				update["$set"].(bson.M)["tags"].(primitive.A),
				mongo.DBRef{Ref: "tags", ID: tag.ID},
			)
		}
	}
	if featuredImage := q.FeaturedImage(); featuredImage != nil {
		if featuredImage.ID.IsZero() {
			update["$set"].(bson.M)["featuredImage"] = mongo.DBRef{}
		} else {
			update["$set"].(bson.M)["featuredImage"] = mongo.DBRef{Ref: "files", ID: featuredImage.ID}
		}
	}
	if attachments := q.Attachments(); attachments != nil {
		update["$set"].(bson.M)["attachments"] = make(primitive.A, 0)
		for _, atm := range *attachments {
			update["$set"].(bson.M)["attachments"] = append(
				update["$set"].(bson.M)["attachments"].(primitive.A),
				mongo.DBRef{Ref: "files", ID: atm.ID},
			)
		}
	}

	_, err := repo.col.UpdateOne(ctx, bson.M{"_id": id.(primitive.ObjectID)}, update)
	if err != nil {
		return Post{}, err
	}

	return repo.FindByID(ctx, id.(primitive.ObjectID))
}

// NewPostQueryBuilder returns a query builder for building post query object
func NewPostQueryBuilder() *PostQueryBuilder {
	return &PostQueryBuilder{postQuery: PostQuery{offset: 0, limit: 5}}
}

// PostQueryBuilder a post object specific query builder
type PostQueryBuilder struct {
	postQuery PostQuery
}

// WithTitle allows to set title to the post query object
func (qb *PostQueryBuilder) WithTitle(title string) *PostQueryBuilder {
	qb.postQuery.title = &title
	return qb
}

// WithSlug allows to set slug to the post query object
func (qb *PostQueryBuilder) WithSlug(slug string) *PostQueryBuilder {
	qb.postQuery.slug = &slug
	return qb
}

// WithStatus allows to set status to the post query object
func (qb *PostQueryBuilder) WithStatus(status Status) *PostQueryBuilder {
	qb.postQuery.status = &status
	return qb
}

// WithMarkdown allows to set markdown to the post query object
func (qb *PostQueryBuilder) WithMarkdown(markdown string) *PostQueryBuilder {
	qb.postQuery.markdown = &markdown
	return qb
}

// WithHTML allows to set HTML to the post query object
func (qb *PostQueryBuilder) WithHTML(html string) *PostQueryBuilder {
	qb.postQuery.html = &html
	return qb
}

// WithPublishedAt allows to set a date-time which the post was published
func (qb *PostQueryBuilder) WithPublishedAt(publishedAt time.Time) *PostQueryBuilder {
	qb.postQuery.publishedAt = &publishedAt
	return qb
}

// WithAuthorID allows to set an author ID to the post query object
func (qb *PostQueryBuilder) WithAuthorID(authorID string) *PostQueryBuilder {
	qb.postQuery.authorID = &authorID
	return qb
}

// WithCategory allows to set category to the post query object
func (qb *PostQueryBuilder) WithCategory(category Category) *PostQueryBuilder {
	qb.postQuery.category = &category
	return qb
}

// WithCategories allows to set list of categories to the post query object
func (qb *PostQueryBuilder) WithCategories(categories []Category) *PostQueryBuilder {
	qb.postQuery.categories = &categories
	return qb
}

// WithTag allows to set tag to the post query object
func (qb *PostQueryBuilder) WithTag(tag Tag) *PostQueryBuilder {
	qb.postQuery.tag = &tag
	return qb
}

// WithTags allows to set list of tags to the post query object
func (qb *PostQueryBuilder) WithTags(tags []Tag) *PostQueryBuilder {
	qb.postQuery.tags = &tags
	return qb
}

// WithFeaturedImage allows to set featured image to the post query object
func (qb *PostQueryBuilder) WithFeaturedImage(featuredImage storage.File) *PostQueryBuilder {
	qb.postQuery.featuredImage = &featuredImage
	return qb
}

// WithAttachments allows to set list of attachments to the post query object
func (qb *PostQueryBuilder) WithAttachments(attachments []storage.File) *PostQueryBuilder {
	qb.postQuery.attachments = &attachments
	return qb
}

// WithOffset allows to set offset to the post query object
func (qb *PostQueryBuilder) WithOffset(offset int64) *PostQueryBuilder {
	qb.postQuery.offset = offset
	return qb
}

// WithLimit allows to set limit the post query object
func (qb *PostQueryBuilder) WithLimit(limit int64) *PostQueryBuilder {
	qb.postQuery.limit = limit
	return qb
}

// Build returns a post query object
func (qb *PostQueryBuilder) Build() PostQuery {
	return qb.postQuery
}

// PostQuery uses as medium for communicating between repository and data-access object (DAO)
type PostQuery struct {
	title         *string
	slug          *string
	status        *Status
	markdown      *string
	html          *string
	publishedAt   *time.Time
	authorID      *string
	category      *Category
	categories    *[]Category
	tag           *Tag
	tags          *[]Tag
	featuredImage *storage.File
	attachments   *[]storage.File

	offset int64
	limit  int64
}

// Title returns title value
func (q PostQuery) Title() *string {
	return q.title
}

// Slug returns slug value
func (q PostQuery) Slug() *string {
	return q.slug
}

// Status return status value
func (q PostQuery) Status() *Status {
	return q.status
}

// Markdown returns markdown value
func (q PostQuery) Markdown() *string {
	return q.markdown
}

// HTML returns HTML value
func (q PostQuery) HTML() *string {
	return q.html
}

// PublishedAt returns date-time value
func (q PostQuery) PublishedAt() *time.Time {
	return q.publishedAt
}

// AuthorID returns author ID
func (q PostQuery) AuthorID() *string {
	return q.authorID
}

// Category returns category object
func (q PostQuery) Category() *Category {
	return q.category
}

// Categories returns list of categories
func (q PostQuery) Categories() *[]Category {
	return q.categories
}

// Tag returns tag object
func (q PostQuery) Tag() *Tag {
	return q.tag
}

// Tags return list of tags
func (q PostQuery) Tags() *[]Tag {
	return q.tags
}

// FeaturedImage returns file object
func (q PostQuery) FeaturedImage() *storage.File {
	return q.featuredImage
}

// Attachments returns list of files object
func (q PostQuery) Attachments() *[]storage.File {
	return q.attachments
}

// Offset returns offset value
func (q PostQuery) Offset() int64 {
	return q.offset
}

// Limit returns limit value
func (q PostQuery) Limit() int64 {
	return q.limit
}
