package storage

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDiskCache_Close(t *testing.T) {
	// Given
	c, _ := NewDiskCache(afero.NewMemMapFs(), ".cache")

	// When
	c.Close()

	// Then
	assert.Panics(t, func() {
		c.doneCh <- struct{}{}
	})
}

func TestDiskCache_Delete(t *testing.T) {
	// Given
	c, _ := NewDiskCache(afero.NewMemMapFs(), ".cache")
	_ = c.Store(bytes.NewReader([]byte("test")), "test")

	// When
	err := c.Delete("test")

	// Then
	assert.Nil(t, err)
	assert.False(t, c.Exists("test"))
}

func TestDiskCache_Exists(t *testing.T) {
	// Given
	c, _ := NewDiskCache(afero.NewMemMapFs(), ".cache")
	_ = c.Store(bytes.NewReader([]byte("test")), "test")

	// When
	result := c.Exists("test")

	// Then
	assert.True(t, result)
}

func TestDiskCache_Retrieve(t *testing.T) {
	c, _ := NewDiskCache(afero.NewMemMapFs(), ".cache")
	_ = c.Store(bytes.NewReader([]byte("test")), "test")

	t.Run("With successful retrieving cache file", func(t *testing.T) {
		// Given

		// When
		body, err := c.Retrieve("test")

		// Then
		assert.Nil(t, err)
		val, _ := ioutil.ReadAll(body)
		assert.Equal(t, "test", string(val))
	})

	t.Run("When unable to retrieving cache file", func(t *testing.T) {
		// Given

		// When
		_, err := c.Retrieve("test2")

		// Then
		assert.NotNil(t, err)
	})
}

func TestDiskCache_Store(t *testing.T) {
	workingDirectory, _ := os.Getwd()
	fs := afero.NewMemMapFs()
	c, _ := NewDiskCache(fs, filepath.Join(workingDirectory, ".cache"))

	t.Run("With successful storing cache file", func(t *testing.T) {
		// Given

		// When
		err := c.Store(bytes.NewReader([]byte("test")), "test")

		// Then
		assert.Nil(t, err)
		assert.True(t, c.Exists("test"))
	})

	t.Run("When unable to make all directories", func(t *testing.T) {
		//// Given
		//
		//// When
		//err := c.Store(bytes.NewReader([]byte("test")), "test")
		//
		//// Then
		//assert.EqualError(t, err, "")
	})
}
