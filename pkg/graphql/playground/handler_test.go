package playground_test

import (
	"github.com/nomkhonwaan/myblog/pkg/graphql/playground"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// Given
	recorder := httptest.NewRecorder()
	data := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut pretium sed ex eget porta. Curabitur placerat condimentum dapibus. In suscipit, massa sed posuere convallis, nisl augue tempor lacus, in semper massa nisi ac lectus. Integer volutpat luctus neque id cursus. Proin ornare rhoncus risus vel elementum. Sed fermentum nulla vel augue venenatis, in vulputate ligula facilisis. Integer scelerisque condimentum ex cursus sagittis. Ut imperdiet ante et scelerisque tempor. In ultrices auctor ex, vitae feugiat nisi interdum et. Nullam efficitur iaculis dolor, eu maximus purus molestie quis. Proin sagittis dui quis iaculis suscipit. Donec id erat aliquam, suscipit enim ac, egestas lectus. Vestibulum pellentesque, justo ac porta facilisis, ex risus tristique erat, ac pharetra dui arcu in elit.")

	// When
	playground.Handler(data).ServeHTTP(recorder, nil)

	// Then
	assert.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))
	assert.Equal(t, "text/html", recorder.Header().Get("Content-Type"))
	assert.Equal(t, data, recorder.Body.Bytes())
}
