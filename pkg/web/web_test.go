package web_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/web"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSPAHandler_ServeHTTP(t *testing.T) {
	var (
		staticFilesPath = os.TempDir()
		h               = NewSPAHandler(staticFilesPath)
	)

	err := ioutil.WriteFile(filepath.Join(staticFilesPath, "index.html"), []byte("test"), 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = os.Remove(filepath.Join(staticFilesPath, "index.html"))
	}()

	newRequest := func(path string) *http.Request {
		return httptest.NewRequest(http.MethodGet, path, nil)
	}

	t.Run("When request to an index path (/)", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		h.ServeHTTP(w, newRequest("/"))

		// Then
		assert.Equal(t, "test", w.Body.String())
	})

	t.Run("When request to non-existing path", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		h.ServeHTTP(w, newRequest("/2019/12/4/test-1"))

		// Then
		assert.Equal(t, "test", w.Body.String())
	})
}
