package graphql

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestServeGraphiqlHandlerFunc(t *testing.T) {
	// Given
	w := httptest.NewRecorder()

	// When
	ServeGraphiqlHandlerFunc([]byte{}).ServeHTTP(w, nil)

	// Then
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
}
