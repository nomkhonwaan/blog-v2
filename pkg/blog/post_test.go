package blog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
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

func TestDecodeFindAllPublishedPostsRequest(t *testing.T) {
	// Given
	defaultPagingOptionsRequest, _ := http.NewRequest("GET", "http://localhost:8080/v1/posts", nil)
	specificPagingOptionsRequest, _ := http.NewRequest("GET", "http://localhost:8080/v1/posts", nil)
	queryParams := url.Values{}
	queryParams.Set("offset", "9")
	queryParams.Set("limit", "10")
	specificPagingOptionsRequest.URL.RawQuery = queryParams.Encode()

	tests := map[string]struct {
		r        *http.Request
		expected findAllPublishedPosts
	}{
		"With default paging options": {
			r:        defaultPagingOptionsRequest,
			expected: findAllPublishedPosts{offset: 0, limit: 5},
		},
		"With specific paging options": {
			r:        specificPagingOptionsRequest,
			expected: findAllPublishedPosts{offset: 9, limit: 10},
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := decodeFindAllPublishedPostsRequest(
				context.Background(),
				test.r,
			)

			assert.Nil(t, err)
			assert.Equal(t, req, test.expected)
		})
	}
	
	// Then
}
