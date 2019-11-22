package graphql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	. "github.com/nomkhonwaan/myblog/pkg/graphql"
	mock_http "github.com/nomkhonwaan/myblog/pkg/http/mock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	slugify "github.com/nomkhonwaan/myblog/pkg/slug"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_RegisterMutation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		catRepo   = mock_blog.NewMockCategoryRepository(ctrl)
		fileRepo  = mock_storage.NewMockFileRepository(ctrl)
		postRepo  = mock_blog.NewMockPostRepository(ctrl)
		tagRepo   = mock_blog.NewMockTagRepository(ctrl)
		transport = mock_http.NewMockRoundTripper(ctrl)

		blogSvc = blog.Service{
			CategoryRepository: catRepo,
			PostRepository:     postRepo,
			TagRepository:      tagRepo,
		}
		fbClient, _ = facebook.NewClient("", "", "", blogSvc, fileRepo, transport)
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

	server := NewServer(blogSvc, fbClient, fileRepo)
	h := Handler(server.Schema())

	t.Run("Create a new post", func(t *testing.T) {
		t.Run("With successful creating a new post", func(t *testing.T) {
			// Given
			q := query{Query: `mutation { createPost { slug } }`}
			w := httptest.NewRecorder()

			postRepo.EXPECT().Create(gomock.Any(), "authorizedID").Return(blog.Post{}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			q := query{Query: `mutation { createPost { slug } }`}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "createPost: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Update post title", func(t *testing.T) {
		t.Run("With successful updating post title", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			title := "Test post"
			slug := fmt.Sprintf("%s-%s", slugify.Make(title), id.Hex())
			q := query{
				Query: `mutation { updatePostTitle(slug: $slug, title: $title) { title slug } }`,
				Variables: map[string]interface{}{
					"slug":  "test-post-" + id.Hex(),
					"title": title,
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			postRepo.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithTitle(title).WithSlug(slug).Build()).Return(blog.Post{}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			title := "Test post"
			q := query{
				Query: `mutation { updatePostTitle(slug: $slug, title: $title) { title slug } }`,
				Variables: map[string]interface{}{
					"slug":  "test-post-" + id.Hex(),
					"title": title,
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostTitle: Unauthorized", result["errors"].([]interface{})[0].(string))
		})

		t.Run("With non-existing post", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			title := "Test post"
			q := query{
				Query: `mutation { updatePostTitle(slug: $slug, title: $title) { title slug } }`,
				Variables: map[string]interface{}{
					"slug":  "test-post-" + id.Hex(),
					"title": title,
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{}, errors.New("test non-existing post"))

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostTitle: test non-existing post", result["errors"].([]interface{})[0].(string))
		})

		t.Run("With different author ID and authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			title := "Test post"
			q := query{
				Query: `mutation { updatePostTitle(slug: $slug, title: $title) { title slug } }`,
				Variables: map[string]interface{}{
					"slug":  "test-post-" + id.Hex(),
					"title": title,
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "otherAuthorID"}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostTitle: Forbidden", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Update post content", func(t *testing.T) {
		t.Run("With successful updating post content", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			markdown := "test"
			html := "<p>test</p>\n"
			q := query{
				Query: `mutation { updatePostContent(slug: $slug, markdown: $markdown) { html markdown } }`,
				Variables: map[string]interface{}{
					"slug":     "test-post-" + id.Hex(),
					"markdown": markdown,
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			postRepo.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithMarkdown(markdown).WithHTML(html).Build())

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			markdown := "test"
			q := query{
				Query: `mutation { updatePostContent(slug: $slug, markdown: $markdown) { html markdown } }`,
				Variables: map[string]interface{}{
					"slug":     "test-post-" + id.Hex(),
					"markdown": markdown,
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostContent: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Update post categories", func(t *testing.T) {
		t.Run("With successful updating post categories", func(t *testing.T) {
			id := primitive.NewObjectID()
			catID := primitive.NewObjectID()
			categories := []blog.Category{{ID: catID}}
			q := query{
				Query: `mutation { updatePostCategories(slug: $slug, categorySlugs: $categorySlugs) { categories { slug } } }`,
				Variables: map[string]interface{}{
					"slug":          "test-post-" + id.Hex(),
					"categorySlugs": []string{"test-category-" + catID.Hex()},
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			catRepo.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{catID}).Return(categories, nil).Times(2)
			postRepo.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithCategories(categories).Build()).Return(blog.Post{Categories: []mongo.DBRef{{ID: catID}}}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)

		})

		t.Run("When unable to retrieve list of categories", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			catID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostCategories(slug: $slug, categorySlugs: $categorySlugs) { categories { slug } } }`,
				Variables: map[string]interface{}{
					"slug":          "test-post-" + id.Hex(),
					"categorySlugs": []string{"test-category-" + catID.Hex()},
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			catRepo.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{catID}).Return(nil, errors.New("test unable to retrieve list of categories"))

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostCategories: test unable to retrieve list of categories", result["errors"].([]interface{})[0].(string))
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			catID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostCategories(slug: $slug, categorySlugs: $categorySlugs) { categories { slug } } }`,
				Variables: map[string]interface{}{
					"slug":          "test-post-" + id.Hex(),
					"categorySlugs": []string{"test-category-" + catID.Hex()},
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostCategories: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Update post tags", func(t *testing.T) {
		t.Run("With successful updating post tags", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			tagID := primitive.NewObjectID()
			tags := []blog.Tag{{ID: tagID}}
			q := query{
				Query: `mutation { updatePostTags(slug: $slug, tagSlugs: $tagSlugs) { tags { slug } } }`,
				Variables: map[string]interface{}{
					"slug":     "test-post-" + id.Hex(),
					"tagSlugs": []string{"test-tag-" + tagID.Hex()},
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			tagRepo.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{tagID}).Return(tags, nil).Times(2)
			postRepo.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithTags(tags).Build()).Return(blog.Post{Tags: []mongo.DBRef{{ID: tagID}}}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When unable to retrieve list of tags", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			tagID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostTags(slug: $slug, tagSlugs: $tagSlugs) { tags { slug } } }`,
				Variables: map[string]interface{}{
					"slug":     "test-post-" + id.Hex(),
					"tagSlugs": []string{"test-tag-" + tagID.Hex()},
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			tagRepo.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{tagID}).Return(nil, errors.New("test unable to retrieve list of tags"))

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostTags: test unable to retrieve list of tags", result["errors"].([]interface{})[0].(string))
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			tagID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostTags(slug: $slug, tagSlugs: $tagSlugs) { tags { slug } } }`,
				Variables: map[string]interface{}{
					"slug":     "test-post-" + id.Hex(),
					"tagSlugs": []string{"test-tag-" + tagID.Hex()},
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostTags: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Update post featured image", func(t *testing.T) {
		t.Run("With successful updating post featured image", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			featuredImageID := primitive.NewObjectID()
			file := storage.File{ID: featuredImageID}
			q := query{
				Query: `mutation { updatePostFeaturedImage(slug: $slug, featuredImageSlug: $featuredImageSlug) { featuredImage { slug } } }`,
				Variables: map[string]interface{}{
					"slug":              "test-post-" + id.Hex(),
					"featuredImageSlug": fmt.Sprintf("featured-image-%s.jpg", featuredImageID.Hex()),
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			fileRepo.EXPECT().FindByID(gomock.Any(), featuredImageID).Return(file, nil).Times(2)
			postRepo.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithFeaturedImage(file).Build()).Return(blog.Post{FeaturedImage: mongo.DBRef{ID: featuredImageID}}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When unable to retrieve featured image by ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			featuredImageID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostFeaturedImage(slug: $slug, featuredImageSlug: $featuredImageSlug) { featuredImage { slug } } }`,
				Variables: map[string]interface{}{
					"slug":              "test-post-" + id.Hex(),
					"featuredImageSlug": fmt.Sprintf("featured-image-%s.jpg", featuredImageID.Hex()),
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			fileRepo.EXPECT().FindByID(gomock.Any(), featuredImageID).Return(storage.File{}, errors.New("test unable to retrieve featured image"))

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostFeaturedImage: test unable to retrieve featured image", result["errors"].([]interface{})[0].(string))
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			featuredImageID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostFeaturedImage(slug: $slug, featuredImageSlug: $featuredImageSlug) { featuredImage { slug } } }`,
				Variables: map[string]interface{}{
					"slug":              "test-post-" + id.Hex(),
					"featuredImageSlug": fmt.Sprintf("featured-image-%s.jpg", featuredImageID.Hex()),
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostFeaturedImage: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})

	t.Run("Update post attachments", func(t *testing.T) {
		t.Run("With successful updating post attachments", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			attachmentID := primitive.NewObjectID()
			attachments := []storage.File{{ID: attachmentID}}
			q := query{
				Query: `mutation { updatePostAttachments(slug: $slug, attachmentSlugs: $attachmentSlugs) { attachments { slug } } }`,
				Variables: map[string]interface{}{
					"slug": "test-post-" + id.Hex(),
					"attachmentSlugs": []string{
						fmt.Sprintf("test-image-%s.jpg", attachmentID.Hex()),
					},
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			fileRepo.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{attachmentID}).Return(attachments, nil).Times(2)
			postRepo.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithAttachments(attachments).Build()).Return(blog.Post{Attachments: []mongo.DBRef{{ID: attachmentID}}}, nil)

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			assert.Equal(t, "200 OK", w.Result().Status)
		})

		t.Run("When unable to retrieve list of attachments", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			attachmentID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostAttachments(slug: $slug, attachmentSlugs: $attachmentSlugs) { attachments { slug } } }`,
				Variables: map[string]interface{}{
					"slug": "test-post-" + id.Hex(),
					"attachmentSlugs": []string{
						fmt.Sprintf("test-image-%s.jpg", attachmentID.Hex()),
					},
				},
			}
			w := httptest.NewRecorder()

			postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{AuthorID: "authorizedID"}, nil)
			fileRepo.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{attachmentID}).Return(nil, errors.New("test unable to retrieve list of attachments"))

			// When
			h.ServeHTTP(w, withAuthorizedID(newGraphQLRequest(q)))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostAttachments: test unable to retrieve list of attachments", result["errors"].([]interface{})[0].(string))
		})

		t.Run("When unable to retrieve authorized ID", func(t *testing.T) {
			// Given
			id := primitive.NewObjectID()
			attachmentID := primitive.NewObjectID()
			q := query{
				Query: `mutation { updatePostAttachments(slug: $slug, attachmentSlugs: $attachmentSlugs) { attachments { slug } } }`,
				Variables: map[string]interface{}{
					"slug": "test-post-" + id.Hex(),
					"attachmentSlugs": []string{
						fmt.Sprintf("test-image-%s.jpg", attachmentID.Hex()),
					},
				},
			}
			w := httptest.NewRecorder()

			// When
			h.ServeHTTP(w, newGraphQLRequest(q))

			// Then
			var result map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&result)

			assert.Equal(t, "updatePostAttachments: Unauthorized", result["errors"].([]interface{})[0].(string))
		})
	})
}
