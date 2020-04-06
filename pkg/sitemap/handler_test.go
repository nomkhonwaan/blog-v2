package sitemap

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/nomkhonwaan/myblog/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeSiteMapHandlerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		cache = mock_storage.NewMockCache(ctrl)
	)

	t.Run("With successful serving sitemap.xml", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://localhost</loc><priority>1</priority></url></urlset>`

		cache.EXPECT().Exists("sitemap.xml").Return(false)
		cache.EXPECT().Store(gomock.Any(), "sitemap.xml").Return(nil)

		// When
		ServeSiteMapHandlerFunc(cache, GenerateFixedURLs("http://localhost")).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://localhost/sitemap.xml", nil))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("With existing sitemap.xml on cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`

		cache.EXPECT().Exists("sitemap.xml").Return(true)
		cache.EXPECT().Retrieve("sitemap.xml").Return(ioutil.NopCloser(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`)), nil)

		// When
		ServeSiteMapHandlerFunc(cache).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://localhost/sitemap.xml", nil))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When unable to generate URLs", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		expected := "test unable to generate URLs\n"

		cache.EXPECT().Exists("sitemap.xml").Return(false)

		// When
		ServeSiteMapHandlerFunc(cache, func() ([]URL, error) { return nil, errors.New("test unable to generate URLs") }).
			ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://localhost/sitemap.xml", nil))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When unable to retrieve sitemap.xml from cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`

		cache.EXPECT().Exists("sitemap.xml").Return(true)
		cache.EXPECT().Retrieve("sitemap.xml").Return(nil, errors.New("test unable to retrieve sitemap.xml from cache"))
		cache.EXPECT().Store(gomock.Any(), "sitemap.xml").Return(nil)

		// When
		ServeSiteMapHandlerFunc(cache).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://localhost/sitemap.xml", nil))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When unable to store new sitemap.xml to cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`

		cache.EXPECT().Exists("sitemap.xml").Return(false)
		cache.EXPECT().Store(gomock.Any(), "sitemap.xml").Return(errors.New("test unable to store new sitemap.xml to cache"))

		// When
		ServeSiteMapHandlerFunc(cache).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://localhost/sitemap.xml", nil))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})
}

func TestGenerateFixedURLs(t *testing.T) {
	// Given
	expected := []URL{
		{
			Location: "http://localhost",
			Priority: 1,
		},
	}

	// When
	urls, err := GenerateFixedURLs("http://localhost")()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expected, urls)
}

func TestGeneratePostURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful generating all post URLs", func(t *testing.T) {
		// Given
		now := time.Now()
		expected := []URL{
			{
				Location:   fmt.Sprintf("http://localhost/%s/test-1", now.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2")),
				LastModify: now.Format(time.RFC3339),
				Priority:   0.8,
			},
			{
				Location:   fmt.Sprintf("http://localhost/%s/test-2", now.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2")),
				LastModify: now.Format(time.RFC3339),
				Priority:   0.8,
			},
		}

		repository.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).
			WithLimit(9999).Build()).Return([]blog.Post{{Slug: "test-1", PublishedAt: now}, {Slug: "test-2", PublishedAt: now, UpdatedAt: now}}, nil)

		// When
		urls, err := GeneratePostURLs("http://localhost", repository)()

		// Then
		assert.Nil(t, err)
		assert.Equal(t, expected, urls)
	})

	t.Run("When unable to find all posts", func(t *testing.T) {
		// Given
		repository.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to find all posts"))

		// When
		_, err := GeneratePostURLs("http://localhost", repository)()

		// Then
		assert.EqualError(t, err, "test unable to find all posts")
	})
}

func TestGenerateCategoryURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockCategoryRepository(ctrl)
	)

	t.Run("With successful generating all category URLs", func(t *testing.T) {
		// Given
		expected := []URL{
			{
				Location: "http://localhost/category/test",
				Priority: 0.5,
			},
		}

		repository.EXPECT().FindAll(gomock.Any()).Return([]blog.Category{{Name: "Test", Slug: "test"}}, nil)

		// When
		urls, err := GenerateCategoryURLs("http://localhost", repository)()

		// Then
		assert.Nil(t, err)
		assert.Equal(t, expected, urls)
	})

	t.Run("When unable to find all categories", func(t *testing.T) {
		// Given
		repository.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("test unable to find all categories"))

		// When
		_, err := GenerateCategoryURLs("http://localhost", repository)()

		// Then
		assert.EqualError(t, err, "test unable to find all categories")
	})
}

func TestGenerateTagURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockTagRepository(ctrl)
	)

	t.Run("With successful generating all tag URLs", func(t *testing.T) {
		// Given
		expected := []URL{
			{
				Location: "http://localhost/tag/test",
				Priority: 0.5,
			},
		}

		repository.EXPECT().FindAll(gomock.Any()).Return([]blog.Tag{{Name: "Test", Slug: "test"}}, nil)

		// When
		urls, err := GenerateTagURLs("http://localhost", repository)()

		// Then
		assert.Nil(t, err)
		assert.Equal(t, expected, urls)
	})

	t.Run("When unable to find all tags", func(t *testing.T) {
		// Given
		repository.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("test unable to find all tags"))

		// When
		_, err := GenerateTagURLs("http://localhost", repository)()

		// Then
		assert.EqualError(t, err, "test unable to find all tags")
	})
}
