package storage

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalDiskCache_Exist(t *testing.T) {
	filePath := filepath.Join(os.TempDir(), "myblog")
	cache, _ := NewLocalDiskCache(filePath)

	defer func() {
		_ = os.Remove(filePath)
	}()

	t.Run("With existing cache file", func(t *testing.T) {
		// Given
		path := "test.txt"
		_ = cache.Store(bytes.NewBufferString("test"), path)
		defer func() {
			_ = os.Remove(filepath.Join(filePath, path))
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
			_ = os.Remove(filepath.Join(filePath, path))
		}

		// When
		result := cache.Exist(path)

		// Then
		assert.False(t, result)
	})
}

func TestLocalDiskCache_Retrieve(t *testing.T) {
	filePath := filepath.Join(os.TempDir(), "myblog")
	cache, _ := NewLocalDiskCache(filePath)
	defer func() {
		_ = os.Remove(filePath)
	}()

	t.Run("With successful retrieving cache file", func(t *testing.T) {
		// Given
		path := "test1.txt"
		_ = cache.Store(bytes.NewBufferString("test"), path)
		defer func() {
			_ = os.Remove(filepath.Join(filePath, path))
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

func TestLocalDiskCache_Store(t *testing.T) {
	filePath := filepath.Join(os.TempDir(), "myblog")
	cache, _ := NewLocalDiskCache(filePath)
	defer func() {
		_ = os.Remove(filePath)
	}()

	t.Run("With successful storing cache file", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("test")
		path := "test3.txt"
		defer func() {
			_ = os.Remove(filepath.Join(filePath, path))
		}()

		// When
		err := cache.Store(body, path)

		// Then
		assert.Nil(t, err)
		val, _ := ioutil.ReadFile(filepath.Join(filePath, path))
		assert.Equal(t, "test", string(val))
	})

	t.Run("When unable to create nested directory on the given path", func(t *testing.T) {
		// Given
		path := "ro/rw/test4.txt"
		body := bytes.NewBufferString("test")
		_ = os.Mkdir(filepath.Join(filePath, "ro"), 0400)
		defer func() {
			_ = os.Remove(filepath.Join(filePath, "ro"))
		}()

		// When
		err := cache.Store(body, path)

		// Then
		assert.NotNil(t, err)
	})

	t.Run("When unable to write cache file", func(t *testing.T) {
		// Given
		path := "test5.txt"
		body := bytes.NewBufferString("test")
		_, _ = os.Create(filepath.Join(filePath, path))
		_ = os.Chmod(filepath.Join(filePath, path), 0400)
		defer func() {
			_ = os.Remove(filepath.Join(filePath, path))
		}()

		// When
		err := cache.Store(body, path)

		// Then
		assert.NotNil(t, err)
	})
}
