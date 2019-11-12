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
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	. "github.com/nomkhonwaan/myblog/pkg/graphql"
	mock_http "github.com/nomkhonwaan/myblog/pkg/http/mock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
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

func TestServer_RegisterQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		catRepo   = mock_blog.NewMockCategoryRepository(ctrl)
		fileRepo  = mock_storage.NewMockFileRepository(ctrl)
		postRepo  = mock_blog.NewMockPostRepository(ctrl)
		tagRepo   = mock_blog.NewMockTagRepository(ctrl)
		transport = mock_http.NewMockRoundTripper(ctrl)

		blogService = blog.Service{
			CategoryRepository: catRepo,
			PostRepository:     postRepo,
			TagRepository:      tagRepo,
		}
		fbClient, _ = facebook.NewClient("", "", "", fileRepo, postRepo, transport)
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

	server := NewServer(BlogService{Service: blogService}, fbClient, fileRepo)
	h := Handler(server.Schema())

	t.Run("With successful querying category by its ID", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		cat := blog.Category{ID: id}
		q := query{
			Query: `{ category(slug: $slug) { slug latestPublishedPosts(offset: 0, limit: 5) { slug } } }`,
			Variables: map[string]interface{}{
				"slug": "test-" + id.Hex(),
			},
		}
		w := httptest.NewRecorder()

		catRepo.EXPECT().FindByID(gomock.Any(), id).Return(cat, nil)
		postRepo.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithCategory(cat).WithOffset(0).WithLimit(5).Build()).Return(nil, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("With successful querying list of categories", func(t *testing.T) {
		// Given
		q := query{
			Query: `{ categories { slug } }`,
		}
		w := httptest.NewRecorder()

		catRepo.EXPECT().FindAll(gomock.Any()).Return([]blog.Category{}, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("With successful querying tag by its ID", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		tag := blog.Tag{ID: id}
		q := query{
			Query: `{ tag(slug: $slug) { slug latestPublishedPosts(offset: 0, limit: 5) { slug } } }`,
			Variables: map[string]interface{}{
				"slug": "test-" + id.Hex(),
			},
		}
		w := httptest.NewRecorder()

		tagRepo.EXPECT().FindByID(gomock.Any(), id).Return(tag, nil)
		postRepo.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithTag(tag).WithOffset(0).WithLimit(5).Build()).Return(nil, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("With successful querying list of tags", func(t *testing.T) {
		// Given
		q := query{Query: `{ tags { slug } }`}
		w := httptest.NewRecorder()

		tagRepo.EXPECT().FindAll(gomock.Any()).Return([]blog.Tag{}, nil)

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

	t.Run("Find post by its slug", func(t *testing.T) {
		t.Run("With existing published post", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			catID := primitive.NewObjectID()
			tagID := primitive.NewObjectID()
			slug := "test-published-" + id.Hex()
			post := blog.Post{
				ID:         id,
				Slug:       slug,
				Status:     blog.Published,
				Categories: []mongo.DBRef{{ID: catID}},
				Tags:       []mongo.DBRef{{ID: tagID}},
			}
			q := query{
				Query: `{ post(slug: $slug) { slug } }`,
				Variables: map[string]interface{}{
					"slug": slug,
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(post, nil)

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("With existing draft post", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			slug := "test-draft-" + id.Hex()
			q := query{
				Query: `{ post(slug: $slug) { slug } }`,
				Variables: map[string]interface{}{
					"slug": slug,
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Status: blog.Draft, AuthorID: "authorizedID"}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When an error has occurred while querying post", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			slug := "test-published-" + id.Hex()
			q := query{
				Query: `{ post(slug: $slug) { slug } }`,
				Variables: map[string]interface{}{
					"slug": slug,
				},
			}
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
			id := primitive.NewObjectID()
			slug := "test-draft-" + id.Hex()
			q := query{
				Query: `{ post(slug: $slug) { slug } }`,
				Variables: map[string]interface{}{
					"slug": slug,
				},
			}
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
			id := primitive.NewObjectID()
			slug := "test-draft-" + id.Hex()
			q := query{
				Query: `{ post(slug: $slug) { slug } }`,
				Variables: map[string]interface{}{
					"slug": slug,
				},
			}
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
