package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// LocalDiskCache implements Cache on local disk
type LocalDiskCache struct {
	filePath string
}

// NewLocalDiskCache returns new LocalDiskCache instance
func NewLocalDiskCache(filePath string) (LocalDiskCache, error) {
	return LocalDiskCache{filePath: filePath}, os.MkdirAll(filePath, 0755)
}

func (c LocalDiskCache) Exist(path string) bool {
	_, err := os.Stat(filepath.Join(c.filePath, path))
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c LocalDiskCache) Retrieve(path string) (io.Reader, error) {
	f, err := os.Open(c.filePath + string(filepath.Separator) + path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c LocalDiskCache) Store(body io.Reader, path string) error {
	dir := filepath.Dir(filepath.Join(c.filePath, path))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(c.filePath, path), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, body)
	return err
}

// LocalDiskStorage implements Uploader and Downloader with LocalDiskCache
type LocalDiskStorage LocalDiskCache

func (s LocalDiskStorage) Download(_ context.Context, path string) (io.Reader, error) {
	return LocalDiskCache(s).Retrieve(path)
}

func (s LocalDiskStorage) Upload(_ context.Context, body io.Reader, path string) error {
	return LocalDiskCache(s).Store(body, path)
}
