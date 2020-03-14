//go:generate mockgen -destination=./mock/cache_mock.go github.com/nomkhonwaan/myblog/pkg/storage Cache

package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Cache uses to storing or retrieving files from hidden or inaccessible place
type Cache interface {
	Exists(path string) bool
	Retrieve(path string) (io.Reader, error)
	Store(body io.Reader, path string) error
}

// DiskCache implements Cache interface with local disk storage
type DiskCache struct {
	fs                 afero.Fs
	doneCh             chan struct{}
	filePath           string
	expirationDuration time.Duration
}

// NewDiskCache returns a new DiskCache instance
func NewDiskCache(fs afero.Fs, filePath string) (*DiskCache, error) {
	c := &DiskCache{
		fs:                 fs,
		doneCh:             make(chan struct{}),
		filePath:           filePath,
		expirationDuration: time.Hour * 24,
	}
	go c.deleteExpiredFiles()
	return c, c.fs.MkdirAll(c.filePath, 0755)
}

func (c DiskCache) deleteExpiredFiles() {
	for {
		select {
		case <-c.doneCh:
			return
		case <-time.After(time.Millisecond * 100):
		default:
			err := filepath.Walk(c.filePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					if time.Since(info.ModTime()) > c.expirationDuration {
						logrus.Infof("deleting expired cache file: %s", path)
						return c.Delete(path)
					}
				}
				return nil
			})
			if err != nil {
				logrus.Errorf("deleteExpiredFiles: %s", err)
			}
		}
	}
}

// Close closes the doneCh which will breaking an infinite-loop file deletion function
func (c *DiskCache) Close() {
	close(c.doneCh)
}

// Exists checks if a file or directory exists
func (c DiskCache) Exists(path string) bool {
	exists, _ := afero.Exists(c.fs, filepath.Join(c.filePath, path))
	return exists
}

// Delete deletes a single file or directory from the give path
func (c DiskCache) Delete(path string) error {
	return c.fs.Remove(filepath.Join(c.filePath, path))
}

// Retrieve returns a single file content from the given path
func (c DiskCache) Retrieve(path string) (io.Reader, error) {
	return c.fs.Open(filepath.Join(c.filePath, path))
}

// Store writes file content to the given path
func (c DiskCache) Store(body io.Reader, path string) error {
	fullPath := filepath.Join(c.filePath, path)

	dir := filepath.Dir(fullPath)
	if err := c.fs.MkdirAll(dir, 0775); err != nil {
		return err
	}
	f, err := c.fs.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, body)
	return err
}
