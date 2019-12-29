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
	cacheFilePath := filepath.Join(os.TempDir(), "myblog")
	cacheService, _ := NewLocalDiskCache(cacheFilePath)

	defer func() {
		_ = os.Remove(cacheFilePath)
	}()

	t.Run("With existing cache file", func(t *testing.T) {
		// Given
		path := "test.txt"
		_ = cacheService.Store(bytes.NewBufferString("test"), path)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilePath, path))
		}()

		// When
		result := cacheService.Exist(path)

		// Then
		assert.True(t, result)
	})

	t.Run("With non-existing cache file", func(t *testing.T) {
		// Given
		path := "test2.txt"
		if cacheService.Exist(path) {
			_ = os.Remove(filepath.Join(cacheFilePath, path))
		}

		// When
		result := cacheService.Exist(path)

		// Then
		assert.False(t, result)
	})
}

func TestLocalDiskCache_Delete(t *testing.T) {
	// Given
	cacheFilePath := filepath.Join(os.TempDir(), "myblog")
	cacheService, _ := NewLocalDiskCache(cacheFilePath)
	path := "test1.txt"
	_ = cacheService.Store(bytes.NewBufferString("test"), path)
	defer func() {
		_ = os.Remove(filepath.Join(cacheFilePath, path))
		_ = os.Remove(cacheFilePath)
	}()

	// When
	err := cacheService.Delete(path)

	// Then
	assert.Nil(t, err)
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalDiskCache_Retrieve(t *testing.T) {
	cacheFilePath := filepath.Join(os.TempDir(), "myblog")
	cacheService, _ := NewLocalDiskCache(cacheFilePath)
	defer func() {
		_ = os.Remove(cacheFilePath)
	}()

	t.Run("With successful retrieving cache file", func(t *testing.T) {
		// Given
		path := "test1.txt"
		_ = cacheService.Store(bytes.NewBufferString("test"), path)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilePath, path))
		}()

		// When
		body, err := cacheService.Retrieve(path)

		// Then
		assert.Nil(t, err)
		val, _ := ioutil.ReadAll(body)
		assert.Equal(t, "test", string(val))
	})

	t.Run("When unable to retrieving cache file", func(t *testing.T) {
		// Given
		path := "test2.txt"

		// When
		_, err := cacheService.Retrieve(path)

		// Then
		assert.NotNil(t, err)
	})
}

func TestLocalDiskCache_Store(t *testing.T) {
	cacheFilePath := filepath.Join(os.TempDir(), "myblog")
	cacheService, _ := NewLocalDiskCache(cacheFilePath)
	defer func() {
		_ = os.Remove(cacheFilePath)
	}()

	t.Run("With successful storing cache file", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("test")
		path := "test3.txt"
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilePath, path))
		}()

		// When
		err := cacheService.Store(body, path)

		// Then
		assert.Nil(t, err)
		val, _ := ioutil.ReadFile(filepath.Join(cacheFilePath, path))
		assert.Equal(t, "test", string(val))
	})

	t.Run("When unable to create nested directory on the given path", func(t *testing.T) {
		// Given
		path := "ro/rw/test4.txt"
		body := bytes.NewBufferString("test")
		_ = os.Mkdir(filepath.Join(cacheFilePath, "ro"), 0400)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilePath, "ro"))
		}()

		// When
		err := cacheService.Store(body, path)

		// Then
		assert.NotNil(t, err)
	})

	t.Run("When unable to write cache file", func(t *testing.T) {
		// Given
		path := "test5.txt"
		body := bytes.NewBufferString("test")
		_, _ = os.Create(filepath.Join(cacheFilePath, path))
		_ = os.Chmod(filepath.Join(cacheFilePath, path), 0400)
		defer func() {
			_ = os.Remove(filepath.Join(cacheFilePath, path))
		}()

		// When
		err := cacheService.Store(body, path)

		// Then
		assert.NotNil(t, err)
	})
}
