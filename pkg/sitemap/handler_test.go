package sitemap_test

//import (
//	"bytes"
//	"errors"
//	"github.com/golang/mock/gomock"
//	"github.com/gorilla/mux"
//	"github.com/nomkhonwaan/myblog/pkg/blog"
//	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
//	. "github.com/nomkhonwaan/myblog/pkg/sitemap"
//	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
//	"github.com/stretchr/testify/assert"
//	"io"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//)
//
//func TestHandler_Register(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	var (
//		baseURL = "http://localhost:8080"
//
//		cacheService = mock_storage.NewMockCache(ctrl)
//		post         = mock_blog.NewMockPostRepository(ctrl)
//		category     = mock_blog.NewMockCategoryRepository(ctrl)
//		tag          = mock_blog.NewMockTagRepository(ctrl)
//
//		blogService = blog.Service{
//			CategoryRepository: category,
//			PostRepository:     post,
//			TagRepository:      tag,
//		}
//
//		r = mux.NewRouter()
//	)
//
//	newSitemapRequest := func() *http.Request {
//		return httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
//	}
//
//	NewHandler(baseURL, cacheService, blogService).Register(r.PathPrefix("/sitemap.xml").Subrouter())
//
//	t.Run("With successful generate a new sitemap.xml from list of posts, categories and tags", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		now := time.Now()
//		posts := []blog.Post{{Slug: "post-1", PublishedAt: now}, {Slug: "post-2", PublishedAt: now.Add(time.Hour * 48 * -1), UpdatedAt: now}}
//		categories := []blog.Category{{Slug: "category-1"}}
//		tags := []blog.Tag{{Slug: "tag-1"}}
//		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://localhost:8080</loc><lastmod>` + now.Format(time.RFC3339) + `</lastmod><priority>1</priority></url><url><loc>http://localhost:8080/` + now.Format("2006/1/2") + `/post-1</loc><lastmod>` + posts[0].PublishedAt.Format(time.RFC3339) + `</lastmod><priority>0.8</priority></url><url><loc>http://localhost:8080/` + posts[1].PublishedAt.Format("2006/1/2") + `/post-2</loc><lastmod>` + now.Format(time.RFC3339) + `</lastmod><priority>0.8</priority></url><url><loc>http://localhost:8080/category/category-1</loc><priority>0.5</priority></url><url><loc>http://localhost:8080/tag/tag-1</loc><priority>0.5</priority></url></urlset>`
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(false)
//		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithLimit(9999).Build()).Return(posts, nil)
//		category.EXPECT().FindAll(gomock.Any()).Return(categories, nil)
//		tag.EXPECT().FindAll(gomock.Any()).Return(tags, nil)
//		cacheService.EXPECT().Store(gomock.Any(), CacheFilePath).DoAndReturn(func(body io.Reader, path string) error {
//			data, _ := ioutil.ReadAll(body)
//			assert.Equal(t, expected, string(data))
//
//			return nil
//		})
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, expected, w.Body.String())
//	})
//
//	t.Run("When cache file available", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		expected := []byte("test cached sitemap.xml content")
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(true)
//		cacheService.EXPECT().Retrieve(CacheFilePath).Return(bytes.NewReader(expected), nil)
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, string(expected), w.Body.String())
//	})
//
//	t.Run("When cache file available but unable to retrieve file content", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		now := time.Now()
//		posts := []blog.Post{{Slug: "post-1", PublishedAt: now}}
//		categories := []blog.Category{{Slug: "category-1"}}
//		tags := []blog.Tag{{Slug: "tag-1"}}
//		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://localhost:8080</loc><lastmod>` + now.Format(time.RFC3339) + `</lastmod><priority>1</priority></url><url><loc>http://localhost:8080/` + now.Format("2006/1/2") + `/post-1</loc><lastmod>` + now.Format(time.RFC3339) + `</lastmod><priority>0.8</priority></url><url><loc>http://localhost:8080/category/category-1</loc><priority>0.5</priority></url><url><loc>http://localhost:8080/tag/tag-1</loc><priority>0.5</priority></url></urlset>`
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(true)
//		cacheService.EXPECT().Retrieve(CacheFilePath).Return(nil, errors.New("test unable to retrieve file content"))
//		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithLimit(9999).Build()).Return(posts, nil)
//		category.EXPECT().FindAll(gomock.Any()).Return(categories, nil)
//		tag.EXPECT().FindAll(gomock.Any()).Return(tags, nil)
//		cacheService.EXPECT().Store(gomock.Any(), CacheFilePath).Return(nil)
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, expected, w.Body.String())
//	})
//
//	t.Run("When unable to retrieve list of published posts", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(false)
//		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithLimit(9999).Build()).Return(nil, errors.New("test unable to retrieve list of published posts"))
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//
//	t.Run("When unable to retrieve list of categories", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		now := time.Now()
//		posts := []blog.Post{{Slug: "post-1", PublishedAt: now}}
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(false)
//		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithLimit(9999).Build()).Return(posts, nil)
//		category.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("test unable to retrieve list of categories"))
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//
//	t.Run("When unable to retrieve list of tags", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		now := time.Now()
//		posts := []blog.Post{{Slug: "post-1", PublishedAt: now}}
//		categories := []blog.Category{{Slug: "category-1"}}
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(false)
//		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithLimit(9999).Build()).Return(posts, nil)
//		category.EXPECT().FindAll(gomock.Any()).Return(categories, nil)
//		tag.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("test unable to retrieve list of tags"))
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//
//	t.Run("When unable to store cache file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		now := time.Now()
//		posts := []blog.Post{{Slug: "post-1", PublishedAt: now}}
//		categories := []blog.Category{{Slug: "category-1"}}
//		tags := []blog.Tag{{Slug: "tag-1"}}
//		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://localhost:8080</loc><lastmod>` + now.Format(time.RFC3339) + `</lastmod><priority>1</priority></url><url><loc>http://localhost:8080/` + now.Format("2006/1/2") + `/post-1</loc><lastmod>` + now.Format(time.RFC3339) + `</lastmod><priority>0.8</priority></url><url><loc>http://localhost:8080/category/category-1</loc><priority>0.5</priority></url><url><loc>http://localhost:8080/tag/tag-1</loc><priority>0.5</priority></url></urlset>`
//
//		cacheService.EXPECT().Exist(CacheFilePath).Return(false)
//		post.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithLimit(9999).Build()).Return(posts, nil)
//		category.EXPECT().FindAll(gomock.Any()).Return(categories, nil)
//		tag.EXPECT().FindAll(gomock.Any()).Return(tags, nil)
//		cacheService.EXPECT().Store(gomock.Any(), CacheFilePath).Return(errors.New("test unable to store cache file"))
//
//		// When
//		r.ServeHTTP(w, newSitemapRequest())
//
//		// Then
//		assert.Equal(t, expected, w.Body.String())
//	})
//}
