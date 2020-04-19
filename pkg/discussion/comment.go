//go:generate mockgen -destination=./mock/comment_mock.go github.com/nomkhonwaan/myblog/pkg/discussion CommentRepository

package discussion

import (
	"context"
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Comment is a a written remark expressing an opinion or reaction
type Comment struct {
	// Identifier of the comment
	ID primitive.ObjectID `bson:"_id" json:"id" graphql:"-"`

	// Identifier of the author
	AuthorID string `bson:"authorId" json:"authorId" graphql:"authorId"`

	// Date-time that the comment was added
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`

	// Date-time that the comment was edited
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt" graphql:"updatedAt"`
}

// MarshalJSON is a custom JSON marshaling function of comment entity
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

// A CommentRepository interface
type CommentRepository interface {
	Create(ctx context.Context, authorID string) (Comment, error)
	Save(ctx context.Context, id interface{}, q CommentQuery) (Comment, error)
}

// NewCommentRepository returns a MongoCommentRepository instance
func NewCommentRepository(db mongo.Database) MongoCommentRepository {
	return MongoCommentRepository{col: mongo.NewCollection(db.Collection("comments"))}
}

// MongoCommentRepository implements CommentRepository interface
type MongoCommentRepository struct {
	col mongo.Collection
}

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

// NewCommentQueryBuilder returns a query builder for building comment query object
func NewCommentQueryBuilder() *CommentQueryBuilder {
	return &CommentQueryBuilder{commentQuery: CommentQuery{offset: 0, limit: 5}}
}

// CommentQueryBuilder is a comment object specific query builder
type CommentQueryBuilder struct {
	commentQuery CommentQuery
}

// WithAuthorID allows to set an author ID to the comment query object
func (qb *CommentQueryBuilder) WithAuthorID(authorID string) *CommentQueryBuilder {
	qb.commentQuery.authorID = &authorID
	return qb
}

// WithOffset allows to set offset to the comment query object
func (qb *CommentQueryBuilder) WithOffset(offset int64) *CommentQueryBuilder {
	qb.commentQuery.offset = offset
	return qb
}

// WithLimit allows to set limit to the comment query object
func (qb *CommentQueryBuilder) WithLimit(limit int64) *CommentQueryBuilder {
	qb.commentQuery.limit = limit
	return qb
}

// CommentQuery uses as medium for communicating between repository and data-access object (DAO)
type CommentQuery struct {
	authorID *string

	offset int64
	limit  int64
}

// AuthorID returns author ID
func (q CommentQuery) AuthorID() *string {
	return q.authorID
}

// Offset returns offset value
func (q CommentQuery) Offset() int64 {
	return q.offset
}

// Limit returns limit value
func (q CommentQuery) Limit() int64 {
	return q.limit
}
