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
	"io/ioutil"
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
		category  = mock_blog.NewMockCategoryRepository(ctrl)
		file      = mock_storage.NewMockFileRepository(ctrl)
		post      = mock_blog.NewMockPostRepository(ctrl)
		tag       = mock_blog.NewMockTagRepository(ctrl)
		transport = mock_http.NewMockRoundTripper(ctrl)

		blogService = blog.Service{
			CategoryRepository: category,
			PostRepository:     post,
			TagRepository:      tag,
		}

		fbClient, _ = facebook.NewClient("", "", "", blogService, file, transport)
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

	server := NewServer(blogService, fbClient, file)
	h := Handler(server.Schema())

	t.Run("With successful querying category by its ID", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		cat := blog.Category{ID: id, Slug: "category-" + id.Hex()}
		q := query{
			Query: `{ category(slug: $slug) { slug latestPublishedPosts(offset: 0, limit: 5) { slug } } }`,
			Variables: map[string]interface{}{
				"slug": "category-" + id.Hex(),
			},
		}
		w := httptest.NewRecorder()

		category.EXPECT().FindByID(gomock.Any(), id).Return(cat, nil)
		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().
			WithCategory(cat).
			WithStatus(blog.Published).
			WithOffset(0).
			WithLimit(5).
			Build(),
		).Return(nil, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		var result struct {
			Data struct {
				Category blog.Category `json:"category"`
			} `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)

		assert.Nil(t, err)
		assert.Equal(t, cat.Slug, result.Data.Category.Slug)
	})

	t.Run("With successful querying list of categories", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		cat := blog.Category{ID: id, Slug: "category-" + id.Hex()}
		q := query{
			Query: `{ categories { slug } }`,
		}
		w := httptest.NewRecorder()

		category.EXPECT().FindAll(gomock.Any()).Return([]blog.Category{cat}, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		var result struct {
			Data struct {
				Categories []blog.Category `json:"categories"`
			} `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)

		assert.Nil(t, err)
		assert.Equal(t, cat.Slug, result.Data.Categories[0].Slug)
	})

	t.Run("With successful querying tag by its ID", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		tg := blog.Tag{ID: id, Slug: "tag-" + id.Hex()}
		q := query{
			Query: `{ tag(slug: $slug) { slug latestPublishedPosts(offset: 0, limit: 5) { slug } } }`,
			Variables: map[string]interface{}{
				"slug": "tag-" + id.Hex(),
			},
		}
		w := httptest.NewRecorder()

		tag.EXPECT().FindByID(gomock.Any(), id).Return(tg, nil)
		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().
			WithTag(tg).
			WithStatus(blog.Published).
			WithOffset(0).
			WithLimit(5).
			Build(),
		).Return(nil, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		var result struct {
			Data struct {
				Tag blog.Tag `json:"tag"`
			} `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)

		assert.Nil(t, err)
		assert.Equal(t, tg.Slug, result.Data.Tag.Slug)
	})

	t.Run("With successful querying list of tags", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		tg := blog.Tag{ID: id, Slug: "tag-" + id.Hex()}
		q := query{Query: `{ tags { slug } }`}
		w := httptest.NewRecorder()

		tag.EXPECT().FindAll(gomock.Any()).Return([]blog.Tag{tg}, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		var result struct {
			Data struct {
				Tags []blog.Tag `json:"tags"`
			} `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)

		assert.Nil(t, err)
		assert.Equal(t, tg.Slug, result.Data.Tags[0].Slug)
	})

	t.Run("With successful querying the latest published posts", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		p := blog.Post{ID: id, Slug: "post-" + id.Hex()}
		q := query{Query: `{ latestPublishedPosts(offset: 0, limit: 5) { slug } }`}

		w := httptest.NewRecorder()

		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().
			WithStatus(blog.Published).
			WithOffset(0).
			WithLimit(5).
			Build(),
		).Return([]blog.Post{p}, nil)

		// When
		h.ServeHTTP(w, newGraphQLRequest(q))

		// Then
		var result struct {
			Data struct {
				LatestPublishedPosts []blog.Post `json:"latestPublishedPosts"`
			} `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&result)

		assert.Nil(t, err)
		assert.Equal(t, p.Slug, result.Data.LatestPublishedPosts[0].Slug)
	})

	t.Run("Find my posts", func(t *testing.T) {
		t.Run("With successful finding my posts", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			slug := "post-" + id.Hex()
			p := blog.Post{ID: id, Slug: slug, Status: blog.Draft, AuthorID: "authorizedID"}
			q := query{
				Query: `{ myPosts(offset: $offset, limit: $limit) { slug status } }`,
				Variables: map[string]interface{}{
					"offset": 0,
					"limit":  5,
				},
			}
			w := httptest.NewRecorder()

			post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().
				WithAuthorID("authorizedID").
				WithOffset(0).
				WithLimit(5).
				Build(),
			).Return([]blog.Post{p}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result struct {
				Data struct {
					MyPosts []blog.Post `json:"myPosts"`
				} `json:"data"`
			}
			err := json.NewDecoder(w.Body).Decode(&result)

			assert.Nil(t, err)
			assert.Equal(t, p.Slug, result.Data.MyPosts[0].Slug)
			assert.Equal(t, blog.Draft, result.Data.MyPosts[0].Status)
		})

		t.Run("When an error has occurred while querying my posts", func(t *testing.T) {
			// Given
			q := query{
				Query: `{ myPosts(offset: $offset, limit: $limit) { slug status } }`,
				Variables: map[string]interface{}{
					"offset": 0,
					"limit":  5,
				},
			}
			w := httptest.NewRecorder()

			post.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("test find my posts error"))

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "myPosts: test find my posts error", result["errors"].([]interface{})[0].(string))
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			q := query{
				Query: `{ myPosts(offset: $offset, limit: $limit) { title slug createdAt updatedAt } }`,
				Variables: map[string]interface{}{
					"offset": 0,
					"limit":  6,
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "myPosts: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Find post by its slug", func(t *testing.T) {
		t.Run("With existing published post", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			catID := primitive.NewObjectID()
			tagID := primitive.NewObjectID()
			slug := "test-published-" + id.Hex()
			p := blog.Post{
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

			post.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result struct {
				Data struct {
					Post blog.Post `json:"post"`
				} `json:"data"`
			}
			err := json.NewDecoder(w.Body).Decode(&result)

			assert.Nil(t, err)
			assert.Equal(t, p.Slug, result.Data.Post.Slug)
		})

		t.Run("With existing published post and has engagement field", func(t *testing.T) {
			t.Run("With successful getting URL result from the Facebook Graph API", func(t *testing.T) {
				// Given
				id := primitive.NewObjectID()
				catID := primitive.NewObjectID()
				tagID := primitive.NewObjectID()
				slug := "post-" + id.Hex()
				p := blog.Post{
					ID:         id,
					Slug:       slug,
					Status:     blog.Published,
					Categories: []mongo.DBRef{{ID: catID}},
					Tags:       []mongo.DBRef{{ID: tagID}},
				}
				q := query{
					Query: `{ post(slug: $slug) { slug engagement { shareCount } } }`,
					Variables: map[string]interface{}{
						"slug": slug,
					},
				}
				w := httptest.NewRecorder()

				post.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)
				transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(_ *http.Request) (*http.Response, error) {
					return &http.Response{
						Body: ioutil.NopCloser(bytes.NewBufferString(`{"engagement":{"comment_count":1,"comment_plugin_count":2,"reaction_count":3,"share_count":4}}`)),
					}, nil
				})

				// When
				h.ServeHTTP(w, newGraphQLRequest(q))

				// Then
				var res struct {
					Data struct {
						Post struct {
							Engagement struct {
								ShareCount int `json:"shareCount"`
							} `json:"engagement"`
						} `json:"post"`
					} `json:"data"`
				}
				err := json.NewDecoder(w.Body).Decode(&res)

				assert.Nil(t, err)
				assert.Equal(t, 4, res.Data.Post.Engagement.ShareCount)
			})

			t.Run("When unable to connect to the Facebook Graph API", func(t *testing.T) {
				// Given
				id := primitive.NewObjectID()
				catID := primitive.NewObjectID()
				tagID := primitive.NewObjectID()
				slug := "post-" + id.Hex()
				p := blog.Post{
					ID:         id,
					Slug:       slug,
					Status:     blog.Published,
					Categories: []mongo.DBRef{{ID: catID}},
					Tags:       []mongo.DBRef{{ID: tagID}},
				}
				q := query{
					Query: `{ post(slug: $slug) { slug engagement { shareCount } } }`,
					Variables: map[string]interface{}{
						"slug": slug,
					},
				}
				w := httptest.NewRecorder()

				post.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)
				transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(_ *http.Request) (*http.Response, error) {
					return nil, errors.New("unable to connect to Facebook Graph API")
				})

				// When
				h.ServeHTTP(w, newGraphQLRequest(q))

				// Then
				var res struct {
					Data struct {
						Post struct {
							Engagement struct {
								ShareCount int `json:"shareCount"`
							} `json:"engagement"`
						} `json:"post"`
					} `json:"data"`
				}
				err := json.NewDecoder(w.Body).Decode(&res)

				assert.Nil(t, err)
				assert.Equal(t, 0, res.Data.Post.Engagement.ShareCount)
			})
		})

		t.Run("With existing draft post", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			slug := "post-" + id.Hex()
			p := blog.Post{ID: id, Slug: slug, Status: blog.Draft, AuthorID: "authorizedID"}
			q := query{
				Query: `{ post(slug: $slug) { slug } }`,
				Variables: map[string]interface{}{
					"slug": slug,
				},
			}
			w := httptest.NewRecorder()

			post.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(p, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result struct {
				Data struct {
					Post blog.Post `json:"post"`
				} `json:"data"`
			}
			err := json.NewDecoder(w.Body).Decode(&result)

			assert.Nil(t, err)
			assert.Equal(t, p.Slug, result.Data.Post.Slug)
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

			post.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test error on finding post by ID"))

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

			post.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Status: blog.Draft}, nil)

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

			post.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Status: blog.Draft, AuthorID: "otherAuthorizedID"}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "post: Forbidden", result["errors"].([]interface{})[0].(string))
		})
	})

}
