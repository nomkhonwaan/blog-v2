package storage_test

import (
	"bytes"
	. "github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDiskCache_Exist(t *testing.T) {
	cacheFilesPath := filepath.Join(os.TempDir(), "myblog")
	cache, _ := NewDiskCache(cacheFilesPath)
	defer func() {
		_ = os.Remove(cacheFilesPath)
	}()

	t.Run("With existing cache file", func(t *testing.T) {
		// Given
		path := "test.txt"
		_ = cache.Store(bytes.NewBufferString("test"), path)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilesPath, path))
		}()

		// When
		result := cache.Exist(path)

		// Then
		assert.True(t, result)
	})

	t.Run("With non-existing cache file", func(t *testing.T) {
		// Given
		path := "test2.txt"
		if cache.Exist(path) {
			_ = os.Remove(filepath.Join(cacheFilesPath, path))
		}

		// When
		result := cache.Exist(path)

		// Then
		assert.False(t, result)
	})
}

func TestDiskCache_Retrieve(t *testing.T) {
	cacheFilesPath := filepath.Join(os.TempDir(), "myblog")
	cache, _ := NewDiskCache(cacheFilesPath)
	defer func() {
		_ = os.Remove(cacheFilesPath)
	}()

	t.Run("With successful retrieving cache file", func(t *testing.T) {
		// Given
		path := "test.txt"
		_ = cache.Store(bytes.NewBufferString("test"), path)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilesPath, path))
		}()

		// When
		body, err := cache.Retrieve(path)

		// Then
		assert.Nil(t, err)
		val, _ := ioutil.ReadAll(body)
		assert.Equal(t, "test", string(val))
	})

	t.Run("When unable to retrieving cache file", func(t *testing.T) {
		// Given
		path := "test2.txt"

		// When
		_, err := cache.Retrieve(path)

		// Then
		assert.NotNil(t, err)
	})
}

func TestDiskCache_Store(t *testing.T) {
	cacheFilesPath := filepath.Join(os.TempDir(), "myblog")
	cache, _ := NewDiskCache(cacheFilesPath)
	defer func() {
		_ = os.Remove(cacheFilesPath)
	}()

	t.Run("With successful storing cache file", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("test")
		path := "test.txt"
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilesPath, path))
		}()

		// When
		err := cache.Store(body, path)

		// Then
		assert.Nil(t, err)
		val, _ := ioutil.ReadFile(filepath.Join(cacheFilesPath, path))
		assert.Equal(t, "test", string(val))
	})

	t.Run("When unable to create nested directory on the given path", func(t *testing.T) {
		// Given
		path := "ro/rw/test.txt"
		body := bytes.NewBufferString("test")
		_ = os.Mkdir(filepath.Join(cacheFilesPath, "ro"), 0400)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilesPath, "ro"))
		}()

		// When
		err := cache.Store(body, path)

		// Then
		assert.NotNil(t, err)
	})

	t.Run("When unable to write cache file", func(t *testing.T) {
		// Given
		path := "test.txt"
		body := bytes.NewBufferString("test")
		_, _ = os.Create(filepath.Join(cacheFilesPath, path))
		_ = os.Chmod(filepath.Join(cacheFilesPath, path), 0400)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilesPath, path))
		}()

		// When
		err := cache.Store(body, path)

		// Then
		assert.NotNil(t, err)
	})
}
