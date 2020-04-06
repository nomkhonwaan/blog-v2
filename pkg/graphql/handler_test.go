package graphql

import (
	"github.com/golang/mock/gomock"
	mock_http "github.com/nomkhonwaan/myblog/internal/http/mock"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		categoryRepository = mock_blog.NewMockCategoryRepository(ctrl)
		tagRepository      = mock_blog.NewMockTagRepository(ctrl)
		postRepository     = mock_blog.NewMockPostRepository(ctrl)
		fileRepository     = mock_storage.NewMockFileRepository(ctrl)
		transport          = mock_http.NewMockRoundTripper(ctrl)
	)

	s, _ := BuildSchema(
		BuildCategorySchema(categoryRepository),
		BuildTagSchema(tagRepository),
		BuildPostSchema(postRepository),
		BuildFileSchema(fileRepository),
		BuildGraphAPISchema("http://localhost", facebook.NewClient("", transport)),
	)

	// When
	Handler(s)

	// Then
}

func TestServeGraphiqlHandlerFunc(t *testing.T) {
	// Given
	w := httptest.NewRecorder()

	// When
	ServeGraphiqlHandlerFunc([]byte{}).ServeHTTP(w, nil)

	// Then
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
}
