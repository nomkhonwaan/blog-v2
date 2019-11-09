package storage

import (
	"io"
	"os"
	"path/filepath"
)

// Cache uses to keeping or retrieving the uploaded files from the cache storage
type Cache interface {
	// Exist checks the existence of the file
	Exist(path string) bool

	// Retrieve returns file content from the given path
	Retrieve(path string) (io.Reader, error)

	// Store keeps file content in the cache storage
	Store(body io.Reader, path string) error
}

// DiskCache uses the hard-disk drive as a cache storage
type DiskCache struct {
	cacheFilesPath string
}

// NewDiskCache returns new disk storage cache instance
func NewDiskCache(cacheFilesPath string) (DiskCache, error) {
	c := DiskCache{
		cacheFilesPath: cacheFilesPath,
	}
	if !c.Exist(cacheFilesPath) {
		if err := os.MkdirAll(cacheFilesPath, 0755); err != nil {
			return DiskCache{}, err
		}
	}
	return c, nil
}

func (c DiskCache) Exist(path string) bool {
	_, err := os.Stat(c.cacheFilesPath + string(filepath.Separator) + path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c DiskCache) Retrieve(path string) (io.Reader, error) {
	f, err := os.Open(c.cacheFilesPath + string(filepath.Separator) + path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c DiskCache) Store(body io.Reader, path string) error {
	dir := filepath.Dir(c.cacheFilesPath + string(filepath.Separator) + path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(c.cacheFilesPath+string(filepath.Separator)+path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, body)
	return err
}
