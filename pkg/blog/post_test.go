package blog

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPost_MarshalJSON(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	createdAt := time.Now()
	post := Post{
		ID:              id,
		Title:           "Children of Dune",
		Slug:            "children-of-dune-" + id.Hex(),
		Status:          Draft,
		Markdown:        "Integer tincidunt ante vel ipsum. Praesent blandit lacinia erat. Vestibulum sed magna at nunc commodo placerat. Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem.",
		HTML:            "Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.",
		PublishedAt:     time.Time{},
		AuthorID:        "github|c7834cb0-2b79-4d27-a817-520a6420c11b",
		DBRefCategories: []mongo.DBRef{},
		DBRefTags:       []mongo.DBRef{},
		CreatedAt:       createdAt,
		UpdatedAt:       time.Time{},
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

func TestPostQueryBuilder_Build(t *testing.T) {
	// Given
	tests := map[string]struct {
		qb       PostQueryBuilder
		expected PostQuery
	}{
		"With default query builder": {
			qb:       NewPostQueryBuilder(),
			expected: &postQuery{offset: 0, limit: 5},
		},
		"With specific status query builder": {
			qb:       NewPostQueryBuilder().WithStatus(Draft),
			expected: &postQuery{status: Draft, offset: 0, limit: 5},
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
	}{
		"With default query options": {
			q:      &postQuery{},
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(0)
			},
		},
		"With specified offset and limit": {
			q:      &postQuery{offset: 10, limit: 5},
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(10).SetLimit(5)
			},
		},
		"With status draft": {
			q:      &postQuery{status: Draft},
			filter: bson.M{"status": Draft},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(0)
			},
		},
		"With status published": {
			q:      &postQuery{status: Published},
			filter: bson.M{"status": Published},
			options: func() *options.FindOptions {
				options := (&options.FindOptions{}).SetSkip(0).SetLimit(0)
				options.Sort = map[string]interface{}{
					"publishedAt": -1,
				}
				return options
			},
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			col.EXPECT().Find(ctx, test.filter, test.options()).Return(cur, nil)
			cur.EXPECT().Close(ctx).Return(nil)
			cur.EXPECT().Decode(gomock.Any()).Return(nil)

			_, err := repo.FindAll(ctx, test.q)

			assert.Nil(t, err)
		})
	}

	// Then
}
