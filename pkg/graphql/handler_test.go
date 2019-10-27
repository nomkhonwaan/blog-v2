package graphql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	. "github.com/nomkhonwaan/myblog/pkg/graphql"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type query struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		service  = mock_blog.NewMockService(ctrl)
		catRepo  = mock_blog.NewMockCategoryRepository(ctrl)
		postRepo = mock_blog.NewMockPostRepository(ctrl)
		tagRepo  = mock_blog.NewMockTagRepository(ctrl)

		newGraphQLRequest = func(q query) *http.Request {
			v, _ := json.Marshal(q)
			return httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(v))
		}
	)

	service.EXPECT().Category().Return(catRepo).AnyTimes()
	service.EXPECT().Post().Return(postRepo).AnyTimes()
	service.EXPECT().Tag().Return(tagRepo).AnyTimes()

	server := NewServer(service)
	h := Handler(server.Schema())

	t.Run("With successful querying list of categories", func(t *testing.T) {
		// Given
		q := query{Query: `{ categories { slug } }`}

		w := httptest.NewRecorder()

		catRepo.EXPECT().FindAll(gomock.Any()).Return([]blog.Category{}, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("With successful querying the latest published posts", func(t *testing.T) {
		// Given
		q := query{Query: `{ latestPublishedPosts(offset: 0, limit: 5) { slug } }`}

		w := httptest.NewRecorder()

		postRepo.EXPECT().FindAll(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, q blog.PostQuery) ([]blog.Post, error) {
			assert.Equal(t, blog.Published, q.Status())
			assert.EqualValues(t, 0, q.Offset())
			assert.EqualValues(t, 5, q.Limit())

			return make([]blog.Post, 0), nil
		})

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})
}
