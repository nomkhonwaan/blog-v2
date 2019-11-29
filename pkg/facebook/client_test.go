package facebook_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	. "github.com/nomkhonwaan/myblog/pkg/facebook"
	mock_http "github.com/nomkhonwaan/myblog/pkg/http/mock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
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

func TestNewClient(t *testing.T) {
	t.Run("When unable to parse open graph template", func(t *testing.T) {
		// Given
		invalidOpenGraphTemplate := "{{.URL}"

		// When
		_, err := NewClient("", "", invalidOpenGraphTemplate, blog.Service{}, nil, http.DefaultTransport)

		// Then
		assert.EqualError(t, err, "template: facebook-open-graph-template:1: unexpected \"}\" in operand")
	})
}

func TestClient_CrawlerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		file        = mock_storage.NewMockFileRepository(ctrl)
		post        = mock_blog.NewMockPostRepository(ctrl)
		blogService = blog.Service{PostRepository: post}

		baseURL     = "http://localhost:8080"
		ogTemplate  = `{"url":"{{.URL}}","title":"{{.Title}}","description":"{{.Description}}","featuredImage":"{{.FeaturedImage}}"}`
		client, _   = NewClient(baseURL, "", ogTemplate, blogService, file, http.DefaultTransport)
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("OK"))
		})
	)

	h := client.CrawlerHandler(nextHandler)

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
		p := blog.Post{
			ID:            id,
			Title:         "Test post",
			Slug:          "test-post-" + id.Hex(),
			PublishedAt:   now,
			Status:        blog.Published,
			FeaturedImage: mongo.DBRef{ID: featuredImageID},
			Markdown: `this should be a post description
but not this line`,
		}
		f := storage.File{
			Slug: fmt.Sprintf("test-featured-image-%s.jpg", featuredImageID.Hex()),
		}
		w := httptest.NewRecorder()

		post.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)
		file.EXPECT().FindByID(gomock.Any(), featuredImageID).Return(f, nil)

		expected := renderedOpenGraphTemplate{
			URL:           fmt.Sprintf("%s/%s/test-post-%s", baseURL, now.Format("2006/1/2"), id.Hex()),
			Title:         "Test post",
			Description:   "this should be a post description",
			FeaturedImage: fmt.Sprintf("%s/api/v2.1/storage/test-featured-image-%s.jpg", baseURL, featuredImageID.Hex()),
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

		post.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{}, errors.New("test unable to find post by ID"))

		// When
		h.ServeHTTP(w, withFacebookUserAgent(newSinglePageRequest("test-post-"+id.Hex())))

		// Then
		assert.Equal(t, "404 Not Found", w.Result().Status)
	})

	t.Run("With non-published status on the post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		p := blog.Post{
			ID:     id,
			Status: blog.Draft,
		}
		w := httptest.NewRecorder()

		post.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)

		// When
		h.ServeHTTP(w, withFacebookUserAgent(newSinglePageRequest("test-post-"+id.Hex())))

		// Then
		assert.Equal(t, "403 Forbidden", w.Result().Status)
	})
}

func TestClient_GetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		transport = mock_http.NewMockRoundTripper(ctrl)

		baseURL        = "http://localhost:8080"
		appAccessToken = "test-app-access-token"
	)

	client, _ := NewClient(baseURL, appAccessToken, "", blog.Service{}, nil, transport)

	t.Run("With successful getting URL result from the Facebook Graph API", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		expected := URL{}
		expected.Engagement.CommentCount = 1
		expected.Engagement.CommentPluginCount = 2
		expected.Engagement.ReactionCount = 3
		expected.Engagement.ShareCount = 4

		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "/v5.0/", r.URL.Path)
			assert.Equal(t, "access_token=test-app-access-token&fields=engagement&id=http%3A%2F%2Flocalhost%3A8080%2F2006%2F1%2F2%2Ftest-post-"+id.Hex(), r.URL.RawQuery)

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"engagement":{"comment_count":1,"comment_plugin_count":2,"reaction_count":3,"share_count":4}}`)),
			}, nil
		})

		// When
		url, err := client.GetURL("/2006/1/2/test-post-" + id.Hex())

		// Then
		assert.Nil(t, err)
		assert.Equal(t, expected, url)
	})

	t.Run("When unable to connect to the Facebook Graph API", func(t *testing.T) {
		// Given
		transport.EXPECT().RoundTrip(gomock.Any()).Return(nil, errors.New("test unable to connect to Facebook Graph API"))

		// When
		_, err := client.GetURL("")

		// Then
		assert.EqualError(t, err, "test unable to connect to Facebook Graph API")
	})
}
