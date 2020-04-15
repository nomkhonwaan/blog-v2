package web

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServeStaticHandlerFunc(t *testing.T) {
	staticFilePath := os.TempDir()
	err := ioutil.WriteFile(filepath.Join(staticFilePath, "index.html"), []byte("test"), 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.Remove(filepath.Join(staticFilePath, "index.html")) }()

	t.Run("When request to an index path (/)", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		ServeStaticHandlerFunc(staticFilePath).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

		// Then
		assert.Equal(t, "test", w.Body.String())
	})

	t.Run("When request to non-existing path", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		ServeStaticHandlerFunc(staticFilePath).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/2019/12/4/test-1", nil))

		// Then
		assert.Equal(t, "test", w.Body.String())
	})
}
