package storage

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalDiskCache implements Cache on local disk
type LocalDiskCache struct {
	filePath            string
	defaultCacheFileTTL time.Duration
	doneCh              chan struct{}
}

// NewLocalDiskCache returns new LocalDiskCache instance
func NewLocalDiskCache(filePath string) (LocalDiskCache, error) {
	c := LocalDiskCache{
		filePath:            filePath,
		defaultCacheFileTTL: time.Hour * 24 * 7,
	}

	go func() {
		for {
			select {
			case _ = <-c.doneCh:
				return
			default:
				err := filepath.Walk(c.filePath, c.deleteExpiredCacheFiles)
				if err != nil {
					logrus.Errorf("unable to delete some cache file: %s", err)
				}

				time.Sleep(time.Second * 5)
			}
		}
	}()

	return c, os.MkdirAll(filePath, 0755)
}

func (c LocalDiskCache) deleteExpiredCacheFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		if time.Since(info.ModTime()) > c.defaultCacheFileTTL {
			logrus.Infof("deleting too old cache file: %s", path)
			if err = os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}

// Close closes the interval cache file deletion function
func (c LocalDiskCache) Close() {
	c.doneCh <- struct{}{}
	close(c.doneCh)
}

func (c LocalDiskCache) Exist(path string) bool {
	_, err := os.Stat(filepath.Join(c.filePath, path))
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c LocalDiskCache) Delete(path string) error {
	return os.Remove(filepath.Join(c.filePath, path))
}

func (c LocalDiskCache) Retrieve(path string) (io.Reader, error) {
	f, err := os.Open(filepath.Join(c.filePath, path))
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

func (s LocalDiskStorage) Delete(_ context.Context, path string) error {
	return LocalDiskCache(s).Delete(path)
}

func (s LocalDiskStorage) Download(_ context.Context, path string) (io.Reader, error) {
	return LocalDiskCache(s).Retrieve(path)
}

func (s LocalDiskStorage) Upload(_ context.Context, body io.Reader, path string) error {
	return LocalDiskCache(s).Store(body, path)
}
