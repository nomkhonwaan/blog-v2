package facebook_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	. "github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsFacebookCrawlerRequest(t *testing.T) {
	// Given
	tests := map[string]struct {
		data     string
		expected bool
	}{
		"With the first Facebook's user-agent string": {
			data:     "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
			expected: true,
		},
		"With the second Facebook's user-agent string": {
			data:     "facebookexternalhit/1.1",
			expected: true,
		},
		"With non-Facebook's user-agent string": {
			data:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36",
			expected: false,
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsFacebookCrawlerRequest(test.data))
		})
	}

	// Then
}

func TestIsSingle(t *testing.T) {
	// Given
	tests := map[string]struct {
		url      string
		id       string
		expected bool
	}{
		"With single page URL": {
			url:      "/2006/1/2/test-post-id",
			id:       "id",
			expected: true,
		},
		"With JavaScript file URL": {
			url:      "/category/test-category-id",
			expected: false,
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			id, yes := IsSingle(test.url)

			assert.Equal(t, test.id, id)
			assert.Equal(t, test.expected, yes)
		})
	}

	// Then
}

func TestNewCrawlerMiddleware(t *testing.T) {
	t.Run("When unable to parse open graph template", func(t *testing.T) {
		// Given
		invalidOpenGraphTemplate := "{{.URL}"

		// When
		_, err := NewCrawlerMiddleware(invalidOpenGraphTemplate, nil, nil)

		// Then
		assert.EqualError(t, err, "template: facebook-opengraph-template:1: unexpected \"}\" in operand")
	})
}

func TestCrawlerMiddleware_Handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		postRepo = mock_blog.NewMockPostRepository(ctrl)
		fileRepo = mock_storage.NewMockFileRepository(ctrl)

		openGraphTemplate = `{"url":"{{.URL}}","title":"{{.Title}}","description":"{{.Description}}","featuredImage":"{{.FeaturedImage}}"}`
		mw, _             = NewCrawlerMiddleware(openGraphTemplate, postRepo, fileRepo)
		nextHandler       = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("OK"))
		})
	)

	h := mw.Handler(nextHandler)

	newSinglePageRequest := func(slug string) *http.Request {
		return httptest.NewRequest(http.MethodGet, "/2006/1/2/"+slug, nil)
	}

	withFacebookUserAgent := func(r *http.Request) *http.Request {
		r.Header.Set("User-Agent", "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)")
		return r
	}

	type renderedOpenGraphTemplate struct {
		URL           string `json:"url"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		FeaturedImage string `json:"featuredImage"`
	}

	t.Run("With Facebook's user-agent string on a single page", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		featuredImageID := primitive.NewObjectID()
		now := time.Now().In(DefaultTimeZone)
		post := blog.Post{
			ID:            id,
			Title:         "Test post",
			Slug:          "test-post-" + id.Hex(),
			PublishedAt:   now,
			Status:        blog.Published,
			FeaturedImage: mongo.DBRef{ID: featuredImageID},
			Markdown: `this should be a post description
but not this line`,
		}
		file := storage.File{
			Slug: fmt.Sprintf("test-featured-image-%s.jpg", featuredImageID.Hex()),
		}
		w := httptest.NewRecorder()

		postRepo.EXPECT().FindByID(gomock.Any(), id).Return(post, nil)
		fileRepo.EXPECT().FindByID(gomock.Any(), featuredImageID).Return(file, nil)

		expected := renderedOpenGraphTemplate{
			URL:           fmt.Sprintf("https://beta.nomkhonwaan.com/%s/test-post-%s", now.Format("2006/1/2"), id.Hex()),
			Title:         "Test post",
			Description:   "this should be a post description",
			FeaturedImage: fmt.Sprintf("https://beta.nomkhonwaan.com/api/v2/storage/test-featured-image-%s.jpg", featuredImageID.Hex()),
		}

		// When
		h.ServeHTTP(w, withFacebookUserAgent(newSinglePageRequest("test-post-"+id.Hex())))

		// Then
		var res renderedOpenGraphTemplate
		err := json.NewDecoder(w.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("With non-Facebook's user-agent request", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		w := httptest.NewRecorder()

		// When
		h.ServeHTTP(w, newSinglePageRequest("test-post-"+id.Hex()))

		// Then
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("When unable to find post by ID", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		w := httptest.NewRecorder()

		postRepo.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{}, errors.New("test unable to find post by ID"))

		// When
		h.ServeHTTP(w, withFacebookUserAgent(newSinglePageRequest("test-post-"+id.Hex())))

		// Then
		assert.Equal(t, "404 Not Found", w.Result().Status)
	})

	t.Run("With non-published status on the post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		post := blog.Post{
			ID:     id,
			Status: blog.Draft,
		}
		w := httptest.NewRecorder()

		postRepo.EXPECT().FindByID(gomock.Any(), id).Return(post, nil)

		// When
		h.ServeHTTP(w, withFacebookUserAgent(newSinglePageRequest("test-post-"+id.Hex())))

		// Then
		assert.Equal(t, "403 Forbidden", w.Result().Status)
	})
}
