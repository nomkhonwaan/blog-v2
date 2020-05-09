//go:generate mockgen -destination=./mock/comment_mock.go github.com/nomkhonwaan/myblog/pkg/discussion CommentRepository

package discussion

import (
	"context"
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Comment is a a written remark expressing an opinion or reaction.
type Comment struct {
	// Identifier of the comment
	ID primitive.ObjectID `bson:"_id" json:"id" graphql:"-"`

	// A parent comment that the comment replies to, default is nil
	Parent mongo.DBRef `bson:"parent" json:"-" graphql:"-"`

	// Identifier of the author
	AuthorID string `bson:"authorId" json:"authorId" graphql:"authorId"`

	// Comment content in plain text string
	Text string `bson:"text" json:"text" graphql:"text"`

	// List of children (that reply to this comment) are belonging to the comment
	Comments []mongo.DBRef `bson:"children" json:"-" graphql:"-"`

	// Date-time that the comment was added
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`

	// Date-time that the comment was edited
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt" graphql:"updatedAt"`
}

// MarshalJSON is a custom JSON marshaling function of comment entity.
func (c Comment) MarshalJSON() ([]byte, error) {
	type Alias Comment
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    c.ID.Hex(),
		Alias: (*Alias)(&c),
	})
}

// A CommentRepository interface.
type CommentRepository interface {
	Create(ctx context.Context, authorID string) (Comment, error)
	FindAllByIDs(ctx context.Context, ids interface{}) ([]Comment, error)
	FindByID(ctx context.Context, id interface{}) (Comment, error)
	Save(ctx context.Context, id interface{}, q CommentQuery) (Comment, error)
}

// NewCommentRepository returns a MongoCommentRepository instance.
func NewCommentRepository(db mongo.Database) MongoCommentRepository {
	return MongoCommentRepository{col: mongo.NewCollection(db.Collection("children"))}
}

// MongoCommentRepository implements CommentRepository interface.
type MongoCommentRepository struct {
	col mongo.Collection
}

// Create inserts a new empty comment which belongs to the author.
func (repo MongoCommentRepository) Create(ctx context.Context, authorID string) (Comment, error) {
	id := primitive.NewObjectID()
	comment := Comment{
		ID:        id,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}

	doc, _ := bson.Marshal(comment)
	_, err := repo.col.InsertOne(ctx, doc)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

// FindAllByIDs returns list of children from list of IDs.
func (repo MongoCommentRepository) FindAllByIDs(ctx context.Context, ids interface{}) ([]Comment, error) {
	filter := bson.M{"_id": bson.M{"$in": ids.([]primitive.ObjectID)}}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := repo.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var comments []Comment
	err = cur.Decode(&comments)

	return comments, err
}

// FindByID returns a single comment from its ID.
func (repo MongoCommentRepository) FindByID(ctx context.Context, id interface{}) (Comment, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var c Comment
	err := r.Decode(&c)

	return c, err
}

// Save does updating a single comment.
func (repo MongoCommentRepository) Save(ctx context.Context, id interface{}, q CommentQuery) (Comment, error) {
	update := bson.M{"$set": bson.M{
		"updatedAt": time.Now(),
	}}

	if parent := q.Parent(); parent != nil {
		update["$set"].(bson.M)["parent"] = mongo.DBRef{Ref: "comments", ID: parent.ID}
	}
	if text := q.Text(); text != nil {
		update["$set"].(bson.M)["text"] = text
	}
	if children := q.Children(); children != nil {
		update["$set"].(bson.M)["children"] = make(primitive.A, 0)
		for _, c := range *children {
			update["$set"].(bson.M)["children"] = append(
				update["$set"].(bson.M)["children"].(primitive.A),
				mongo.DBRef{Ref: "comments", ID: c.ID},
			)
		}
	}

	_, err := repo.col.UpdateOne(ctx, bson.M{"_id": id.(primitive.ObjectID)}, update)
	if err != nil {
		return Comment{}, err
	}

	return repo.FindByID(ctx, id.(primitive.ObjectID))
}

// NewCommentQueryBuilder returns a query builder for building comment query object.
func NewCommentQueryBuilder() *CommentQueryBuilder {
	return &CommentQueryBuilder{commentQuery: CommentQuery{offset: 0, limit: 5}}
}

// CommentQueryBuilder is a comment object specific query builder.
type CommentQueryBuilder struct {
	commentQuery CommentQuery
}

// WithParent allows to set parent to the query object.
func (qb *CommentQueryBuilder) WithParent(c Comment) *CommentQueryBuilder {
	qb.commentQuery.parent = &c
	return qb
}

// WithAuthorID allows to set an author ID to the comment query object.
func (qb *CommentQueryBuilder) WithAuthorID(authorID string) *CommentQueryBuilder {
	qb.commentQuery.authorID = &authorID
	return qb
}

// WithText allows to set text to the comment query object.
func (qb *CommentQueryBuilder) WithText(text string) *CommentQueryBuilder {
	qb.commentQuery.text = &text
	return qb
}

// WithChildren allows to set list of comments to the comment query object.
func (qb *CommentQueryBuilder) WithChildren(comments []Comment) *CommentQueryBuilder {
	qb.commentQuery.children = &comments
	return qb
}

// WithOffset allows to set offset to the comment query object.
func (qb *CommentQueryBuilder) WithOffset(offset int64) *CommentQueryBuilder {
	qb.commentQuery.offset = offset
	return qb
}

// WithLimit allows to set limit to the comment query object.
func (qb *CommentQueryBuilder) WithLimit(limit int64) *CommentQueryBuilder {
	qb.commentQuery.limit = limit
	return qb
}

// Build returns a comment query object.
func (qb *CommentQueryBuilder) Build() CommentQuery {
	return qb.commentQuery
}

// CommentQuery uses as medium for communicating between repository and data-access object (DAO).
type CommentQuery struct {
	parent   *Comment
	authorID *string
	text     *string
	children *[]Comment

	offset int64
	limit  int64
}

// Parent returns parent comment.
func (q CommentQuery) Parent() *Comment {
	return q.parent
}

// AuthorID returns author ID.
func (q CommentQuery) AuthorID() *string {
	return q.authorID
}

// Text returns text value.
func (q CommentQuery) Text() *string {
	return q.text
}

// Children returns list of comment object.
func (q CommentQuery) Children() *[]Comment {
	return q.children
}

// Offset returns offset value.
func (q CommentQuery) Offset() int64 {
	return q.offset
}

// Limit returns limit value.
func (q CommentQuery) Limit() int64 {
	return q.limit
}
