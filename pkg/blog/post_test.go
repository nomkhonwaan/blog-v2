package blog

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJSONMarshalingPostEntity(t *testing.T) {
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
		Categories:  []Category{},
		Tags:        []Tag{},
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
	assert.Equal(t, recorder.Body.String(), "{\"id\":\""+id.Hex()+"\",\"title\":\"Children of Dune\",\"slug\":\"children-of-dune-"+id.Hex()+"\",\"status\":\"DRAFT\",\"markdown\":\"Integer tincidunt ante vel ipsum. Praesent blandit lacinia erat. Vestibulum sed magna at nunc commodo placerat. Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem.\",\"html\":\"Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.\",\"publishedAt\":\"0001-01-01T00:00:00Z\",\"AuthorID\":\"github|c7834cb0-2b79-4d27-a817-520a6420c11b\",\"Categories\":[],\"Tags\":[],\"createdAt\":\""+createdAt.Format(time.RFC3339Nano)+"\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}\n")
}

func TestPostQueryBuilder(t *testing.T) {
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
			assert.Equal(t, test.qb.Build(), test.expected)
		})
	}

	// Then
}
