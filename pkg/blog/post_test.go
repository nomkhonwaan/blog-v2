package blog

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPost_MarshalJSON(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	createdAt := time.Now()
	post := Post{
		ID:          id,
		Title:       "Children of Dune",
		Slug:        "children-of-dune-" + id.Hex(),
		Status:      Draft,
		Markdown:    "Integer tincidunt ante vel ipsum. Praesent blandit lacinia erat. Vestibulum sed magna at nunc commodo placerat. Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem.",
		HTML:        "Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.",
		PublishedAt: time.Time{},
		AuthorID:    "github|c7834cb0-2b79-4d27-a817-520a6420c11b",
		Categories:  []mongo.DBRef{},
		Tags:        []mongo.DBRef{},
		CreatedAt:   createdAt,
		UpdatedAt:   time.Time{},
	}
	recorder := httptest.NewRecorder()

	// When
	err := encodeResponse(
		context.Background(),
		recorder,
		post,
	)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\""+id.Hex()+"\",\"title\":\"Children of Dune\",\"slug\":\"children-of-dune-"+id.Hex()+"\",\"status\":\"DRAFT\",\"markdown\":\"Integer tincidunt ante vel ipsum. Praesent blandit lacinia erat. Vestibulum sed magna at nunc commodo placerat. Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem.\",\"html\":\"Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.\",\"publishedAt\":\"0001-01-01T00:00:00Z\",\"authorId\":\"github|c7834cb0-2b79-4d27-a817-520a6420c11b\",\"createdAt\":\""+createdAt.Format(time.RFC3339Nano)+"\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}\n", recorder.Body.String())
}

func TestPost_BelongToCategories(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	categoryRepo := NewMockCategoryRepository(ctrl)
	ctx := context.Background()
	categoryID := primitive.NewObjectID()
	categories := []Category{{ID: categoryID}}

	categoryRepo.EXPECT().FindAllByIDs(ctx, []primitive.ObjectID{categoryID}).Return(categories, nil)

	post := Post{Categories: []mongo.DBRef{{ID: categoryID}}}

	// When
	result, err := post.BelongToCategories(categoryRepo).(func(context.Context, Post) ([]Category, error))(ctx, post)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, categories, result)
}

func TestPost_BelongToTags(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagRepo := NewMockTagRepository(ctrl)
	ctx := context.Background()
	tagID := primitive.NewObjectID()
	tags := []Tag{{ID: tagID}}

	tagRepo.EXPECT().FindAllByIDs(ctx, []primitive.ObjectID{tagID}).Return(tags, nil)

	post := Post{Tags: []mongo.DBRef{{ID: tagID}}}

	// When
	result, err := post.BelongToTags(tagRepo).(func(context.Context, Post) ([]Tag, error))(ctx, post)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, tags, result)
}

func TestNewPostRepository(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := mongo.NewMockCollection(ctrl)

	// When
	repo := NewPostRepository(col)

	// Then
	assert.Equal(t, MongoPostRepository{col}, repo)
}

func TestMongoPostRepository_Create(t *testing.T) {
	t.Run("When insert into the collection successfully", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		authorID := "github|303589"

		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, nil)

		postRepo := NewPostRepository(col)

		// When
		result, err := postRepo.Create(ctx, authorID)

		// Then
		assert.Nil(t, err)
		assert.True(t, time.Since(result.CreatedAt) < time.Minute)
		assert.Equal(t, fmt.Sprintf("%s", result.ID.Hex()), result.Slug)
		assert.Equal(t, Draft, result.Status)
		assert.Equal(t, authorID, result.AuthorID)
	})

	t.Run("When insert into the collection un-successfully", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		authorID := "github|303589"

		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, errors.New("something went wrong"))

		postRepo := NewPostRepository(col)

		expected := Post{}

		//result When
		result, err := postRepo.Create(ctx, authorID)

		// Then
		assert.EqualError(t, err, "something went wrong")
		assert.Equal(t, expected, result)
	})
}

func TestMongoPostRepository_FindAll(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cur := mongo.NewMockCursor(ctrl)
	col := mongo.NewMockCollection(ctrl)
	ctx := context.Background()

	repo := NewPostRepository(col)

	tests := map[string]struct {
		q       PostQuery
		filter  interface{}
		options func() *options.FindOptions
		err     error
	}{
		"With default query options": {
			q:      &MongoPostQuery{},
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(0)
			},
		},
		"With specified offset and limit": {
			q:      &MongoPostQuery{offset: 10, limit: 5},
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(10).SetLimit(5)
			},
		},
		"With status draft": {
			q:      &MongoPostQuery{status: Draft},
			filter: bson.M{"status": Draft},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(0)
			},
		},
		"With status published": {
			q:      &MongoPostQuery{status: Published},
			filter: bson.M{"status": Published},
			options: func() *options.FindOptions {
				options := (&options.FindOptions{}).SetSkip(0).SetLimit(0)
				options.Sort = map[string]interface{}{
					"publishedAt": -1,
				}
				return options
			},
		},
		"When an error has occurred on finding the result": {
			q:      &MongoPostQuery{},
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(0)
			},
			err: errors.New("something went wrong"),
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			col.EXPECT().Find(ctx, test.filter, test.options()).Return(cur, test.err)

			if test.err == nil {
				cur.EXPECT().Close(ctx).Return(nil)
				cur.EXPECT().Decode(gomock.Any()).Return(nil)
				_, err := repo.FindAll(ctx, test.q)

				assert.Nil(t, err)
			} else {
				_, err := repo.FindAll(ctx, test.q)

				assert.EqualError(t, err, test.err.Error())
			}
		})
	}

	// Then
}

func TestMongoPostRepository_FindByID(t *testing.T) {
}

func TestNewPostQueryBuilder(t *testing.T) {
	// Given
	expected := &MongoPostQueryBuilder{
		MongoPostQuery: &MongoPostQuery{
			offset: 0,
			limit:  5,
		},
	}

	// When
	qb := NewPostQueryBuilder()

	// Then
	assert.Equal(t, expected, qb)
}

func TestPostQueryBuilder_WithStatus(t *testing.T) {
	// Given
	expected := &MongoPostQueryBuilder{
		MongoPostQuery: &MongoPostQuery{
			status: Published,
			offset: 0,
			limit:  5,
		},
	}

	// When
	qb := NewPostQueryBuilder().WithStatus(Published)

	// Then
	assert.Equal(t, expected, qb)
}

func TestPostQueryBuilder_WithOffset(t *testing.T) {
	// Given
	expected := &MongoPostQueryBuilder{
		MongoPostQuery: &MongoPostQuery{
			offset: 99,
			limit:  5,
		},
	}

	// When
	qb := NewPostQueryBuilder().WithOffset(99)

	// Then
	assert.Equal(t, expected, qb)
}

func TestPostQueryBuilder_WithLimit(t *testing.T) {
	// Given
	expected := &MongoPostQueryBuilder{
		MongoPostQuery: &MongoPostQuery{
			offset: 0,
			limit:  99,
		},
	}

	// When
	qb := NewPostQueryBuilder().WithLimit(99)

	// Then
	assert.Equal(t, expected, qb)
}

func TestPostQueryBuilder_Build(t *testing.T) {
	// Given
	tests := map[string]struct {
		qb       PostQueryBuilder
		expected PostQuery
	}{
		"With default query builder": {
			qb:       NewPostQueryBuilder(),
			expected: &MongoPostQuery{offset: 0, limit: 5},
		},
		"With specific status query builder": {
			qb:       NewPostQueryBuilder().WithStatus(Draft),
			expected: &MongoPostQuery{status: Draft, offset: 0, limit: 5},
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.qb.Build())
		})
	}

	// Then
}

func TestPostQuery_Status(t *testing.T) {

}

func TestPostQuery_Offset(t *testing.T) {

}

func TestPostQuery_Limit(t *testing.T) {

}
