package github_test

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/nomkhonwaan/myblog/pkg/github"
	mock_http "github.com/nomkhonwaan/myblog/pkg/http/mock"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		cacheService = mock_storage.NewMockCache(ctrl)
		transport    = mock_http.NewMockRoundTripper(ctrl)

		r = mux.NewRouter()
	)

	newGistRequest := func(src string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/v2.1/github/gist", nil)
		u := req.URL.Query()
		u.Set("src", src)
		req.URL.RawQuery = u.Encode()
		return req
	}

	NewHandler(cacheService, transport).Register(r.PathPrefix("/api/v2.1/github").Subrouter())

	t.Run("With successful retrieving Gist content", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/gist-id.js?file=test.txt"

		cacheService.EXPECT().Exist(url.QueryEscape(src) + ".json").Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "gist.github.com", r.URL.Host)
			assert.Equal(t, "/nomkhonwaan/gist-id.json", r.URL.Path)
			assert.Equal(t, "test.txt", r.URL.Query().Get("file"))

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			}, nil
		})
		cacheService.EXPECT().Store(gomock.Any(), url.QueryEscape(src)+".json").DoAndReturn(func(body io.Reader, path string) error {
			data, _ := ioutil.ReadAll(body)
			assert.Equal(t, `{"foo":"bar"}`, string(data))

			return nil
		})

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
		assert.Equal(t, `{"foo":"bar"}`, w.Body.String())
	})

	t.Run("When cache file available", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/gist-id.js?file=test.txt"

		cacheService.EXPECT().Exist(url.QueryEscape(src) + ".json").Return(true)
		cacheService.EXPECT().Retrieve(url.QueryEscape(src)+".json").Return(bytes.NewBufferString(`{"foo":"bar"}`), nil)

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
		assert.Equal(t, `{"foo":"bar"}`, w.Body.String())
	})

	t.Run("When cache file available but unable to retrieve file content", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/gist-id.js?file=test.txt"

		cacheService.EXPECT().Exist(url.QueryEscape(src) + ".json").Return(true)
		cacheService.EXPECT().Retrieve(url.QueryEscape(src)+".json").Return(nil, errors.New("test unable to retrieve file content"))
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "gist.github.com", r.URL.Host)
			assert.Equal(t, "/nomkhonwaan/gist-id.json", r.URL.Path)
			assert.Equal(t, "test.txt", r.URL.Query().Get("file"))

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			}, nil
		})
		cacheService.EXPECT().Store(gomock.Any(), url.QueryEscape(src)+".json").DoAndReturn(func(body io.Reader, path string) error {
			data, _ := ioutil.ReadAll(body)
			assert.Equal(t, `{"foo":"bar"}`, string(data))

			return nil
		})

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
		assert.Equal(t, `{"foo":"bar"}`, w.Body.String())
	})

	t.Run("When client try to send different Gist host name to the API", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.malicious.localtest.me/nomkhonwaan/gist-id.js?file=test.txt"

		cacheService.EXPECT().Exist(url.QueryEscape(src) + ".json").Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "gist.github.com", r.URL.Host)
			assert.Equal(t, "/nomkhonwaan/gist-id.json", r.URL.Path)
			assert.Equal(t, "test.txt", r.URL.Query().Get("file"))

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			}, nil
		})
		cacheService.EXPECT().Store(gomock.Any(), url.QueryEscape(src)+".json").Return(nil)

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
	})

	t.Run("With empty `src`", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := ""

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
		assert.Equal(t, "400 Bad Request", w.Result().Status)
	})

	t.Run("When unable to retrieve Gist content", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/gist-id.js?file=test.txt"

		cacheService.EXPECT().Exist(url.QueryEscape(src) + ".json").Return(false)
		transport.EXPECT().RoundTrip(gomock.Any()).Return(nil, errors.New("test unable to retrieve Gist content"))

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
		assert.Equal(t, "test unable to retrieve Gist content\n", w.Body.String())
	})

	t.Run("When unable to store cache file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		src := "https://gist.github.com/nomkhonwaan/gist-id.js?file=test.txt"

		cacheService.EXPECT().Exist(url.QueryEscape(src) + ".json").Return(true)
		cacheService.EXPECT().Retrieve(url.QueryEscape(src)+".json").Return(nil, errors.New("test unable to retrieve file content"))
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "gist.github.com", r.URL.Host)
			assert.Equal(t, "/nomkhonwaan/gist-id.json", r.URL.Path)
			assert.Equal(t, "test.txt", r.URL.Query().Get("file"))

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			}, nil
		})
		cacheService.EXPECT().Store(gomock.Any(), url.QueryEscape(src)+".json").Return(errors.New("test unable to store cache file"))

		// When
		r.ServeHTTP(w, newGistRequest(src))

		// Then
	})
}
