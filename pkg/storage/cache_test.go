package storage

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	mock_afero "github.com/nomkhonwaan/myblog/internal/afero/mock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestDiskCache_Close(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fs = mock_afero.NewMockFs(ctrl)
	)

	fs.EXPECT().MkdirAll(".cache", gomock.Any()).Return(nil)

	c, _ := NewDiskCache(fs, ".cache")

	// When
	c.Close()

	// Then
	assert.Panics(t, func() {
		c.doneCh <- struct{}{}
	})
}

func TestDiskCache_Delete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fs = mock_afero.NewMockFs(ctrl)
	)

	fs.EXPECT().MkdirAll(".cache", gomock.Any()).Return(nil)
	fs.EXPECT().Remove(filepath.Join(".cache", "test")).Return(nil)

	c, _ := NewDiskCache(fs, ".cache")
	defer c.Close()

	// When
	err := c.Delete("test")

	// Then
	assert.Nil(t, err)
}

func TestDiskCache_Exists(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fs = mock_afero.NewMockFs(ctrl)
	)

	fs.EXPECT().MkdirAll(".cache", gomock.Any()).Return(nil)
	fs.EXPECT().Stat(filepath.Join(".cache", "test")).Return(nil, nil)

	c, _ := NewDiskCache(fs, ".cache")
	defer c.Close()

	// When
	result := c.Exists("test")

	// Then
	assert.True(t, result)
}

func TestDiskCache_Retrieve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fs   = mock_afero.NewMockFs(ctrl)
		file = mock_afero.NewMockFile(ctrl)
	)

	fs.EXPECT().MkdirAll(".cache", gomock.Any()).Return(nil)

	c, _ := NewDiskCache(fs, ".cache")
	defer c.Close()

	t.Run("With successful retrieving cache file", func(t *testing.T) {
		// Given
		fs.EXPECT().Open(filepath.Join(".cache", "test")).Return(file, nil)

		// When
		_, err := c.Retrieve("test")

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to retrieving cache file", func(t *testing.T) {
		// Given
		fs.EXPECT().Open(gomock.Any()).Return(nil, errors.New("test unable to retrieve cache file"))

		// When
		_, err := c.Retrieve("test")

		// Then
		assert.EqualError(t, err, "test unable to retrieve cache file")
	})
}

func TestDiskCache_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		fs   = mock_afero.NewMockFs(ctrl)
		file = mock_afero.NewMockFile(ctrl)
	)

	fs.EXPECT().MkdirAll(".cache", gomock.Any()).Return(nil)

	c, _ := NewDiskCache(fs, ".cache")
	defer c.Close()

	t.Run("With successful storing cache file", func(t *testing.T) {
		// Given
		fs.EXPECT().MkdirAll(".cache", gomock.Any()).Return(nil)
		fs.EXPECT().OpenFile(filepath.Join(".cache", "test"), gomock.Any(), gomock.Any()).Return(file, nil)
		file.EXPECT().Write([]byte("test")).Return(len("test"), nil)

		// When
		err := c.Store(bytes.NewReader([]byte("test")), "test")

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to make all directories", func(t *testing.T) {
		// Given
		fs.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(errors.New("test unable to make all directories"))

		// When
		err := c.Store(bytes.NewReader([]byte("test")), "test")

		// Then
		assert.EqualError(t, err, "test unable to make all directories")
	})

	t.Run("When unable to open file for write", func(t *testing.T) {
		// Given
		fs.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(nil)
		fs.EXPECT().OpenFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to open file for write"))

		// When
		err := c.Store(bytes.NewReader([]byte("test")), "test")

		// Then
		assert.EqualError(t, err, "test unable to open file for write")
	})
}
