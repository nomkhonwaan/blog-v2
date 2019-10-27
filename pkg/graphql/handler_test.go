package graphql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	. "github.com/nomkhonwaan/myblog/pkg/graphql"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	)

	newGraphQLRequest := func(q query) *http.Request {
		v, _ := json.Marshal(q)
		return httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(v))
	}

	withAuthorizedID := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(context.Background(), auth.UserProperty, &jwt.Token{
			Claims: jwt.MapClaims{
				"sub": "authorizedID",
			},
		}))
	}

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

	t.Run("Find post by its ID", func(t *testing.T) {
		t.Run("With existing published post", func(t *testing.T) {
			// Given
			expected := primitive.NewObjectID()
			slug := "test-published-" + expected.Hex()
			q := query{Query: `{ post(slug: "` + slug + `") { slug } }`}

			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, id interface{}) (blog.Post, error) {
				assert.IsType(t, primitive.ObjectID{}, id)
				assert.Equal(t, expected.Hex(), id.(primitive.ObjectID).Hex())

				return blog.Post{Status: blog.Published}, nil
			})

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("With existing draft post", func(t *testing.T) {
			// Given
			expected := primitive.NewObjectID()
			slug := "test-draft-" + expected.Hex()
			q := query{Query: `{ post(slug: "` + slug + `") { slug } }`}

			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Status: blog.Draft, AuthorID: "authorizedID"}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When an error has occurred while querying post", func(t *testing.T) {
			// Given
			expected := primitive.NewObjectID()
			slug := "test-published-" + expected.Hex()
			q := query{Query: `{ post(slug: "` + slug + `") { slug } }`}

			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test error on finding post by ID"))

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "post: test error on finding post by ID", result["errors"].([]interface{})[0].(string))
		})

		t.Run("With existing draft post but unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			expected := primitive.NewObjectID()
			slug := "test-draft-" + expected.Hex()
			q := query{Query: `{ post(slug: "` + slug + `") { slug } }`}

			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Status: blog.Draft}, nil)

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "post: Unauthorized", result["errors"].([]interface{})[0].(string))
		})

		t.Run("With existing draft post but not post's author", func(t *testing.T) {
			// Given
			expected := primitive.NewObjectID()
			slug := "test-draft-" + expected.Hex()
			q := query{Query: `{ post(slug: "` + slug + `") { slug } }`}

			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Status: blog.Draft, AuthorID: "otherAuthorizedID"}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "post: Forbidden", result["errors"].([]interface{})[0].(string))
		})
	})
}
