package github

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	mock_http "github.com/nomkhonwaan/myblog/internal/http/mock"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetGistHandlerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		cache     = mock_storage.NewMockCache(ctrl)
		transport = mock_http.NewMockRoundTripper(ctrl)
	)

	newGettingGistRequest := func(src string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/v2.1/github/gist", nil)
		q := req.URL.Query()
		q.Set("src", src)
		req.URL.RawQuery = q.Encode()
		return req
	}

	t.Run("With successful getting GitHub Gist from the given URL", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.js?file=fasthttp.go"

		cache.EXPECT().Exists(url.QueryEscape(src) + ".json").Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "gist.github.com", r.URL.Host)
			assert.Equal(t, "/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.json", r.URL.Path)
			assert.Equal(t, "fasthttp.go", r.URL.Query().Get("file"))

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`)),
			}, nil
		})
		cache.EXPECT().Store(gomock.Any(), url.QueryEscape(src)+".json").Return(nil)

		expected := `{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(src))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("With existing Gist file on cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.js?file=fasthttp.go"

		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(bytes.NewBufferString(`{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`)), nil)

		expected := `{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(src))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When unable to retrieve Gist file from cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.js?file=fasthttp.go"

		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(nil, errors.New("test unable to retrieve Gist file from cache"))
		transport.EXPECT().RoundTrip(gomock.Any()).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`)),
		}, nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		expected := `{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(src))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When src host is not gist.github.com", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.malicious.localtest.me/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.js?file=fasthttp.go"

		cache.EXPECT().Exists(gomock.Any()).Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`)),
		}, nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		expected := `{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(src))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When unable to fetch GitHub Gist", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.malicious.localtest.me/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.js?file=fasthttp.go"

		cache.EXPECT().Exists(gomock.Any()).Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).Return(nil, errors.New("test unable to fetch GitHub Gist"))

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(src))

		// Then
		assert.Equal(t, "Get https://gist.github.com/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.json?file=fasthttp.go: test unable to fetch GitHub Gist\n", w.Body.String())
	})

	t.Run("When src is empty", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(""))

		// Then
		assert.Equal(t, "src value is empty\n", w.Body.String())
	})

	t.Run("When unable to store Gist file to cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/b7526527067b1069d73d3b991be8b93c.js?file=fasthttp.go"

		cache.EXPECT().Exists(gomock.Any()).Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`)),
		}, nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errors.New("test unable to store Gist file to cache"))

		expected := `{"div": "<div id=\"gist101439575\" class=\"gist\"></div>", "stylesheet": "https://github.githubassets.com/assets/gist-embed-31007ea0d3bd9f80540adfbc55afc7bd.css"}`

		// When
		GetGistHandlerFunc(cache, transport).ServeHTTP(w, newGettingGistRequest(src))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})
}
