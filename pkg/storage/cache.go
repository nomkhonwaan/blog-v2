// go:generate mockgen -destination=./mock/cache_mock.go github.com/nomkhonwaan/myblog/storage Cache

package storage

import (
	"io"
)

// Cache keeps data at the nearest place where application can store and retrieve quickly
type Cache interface {
	// Check an existing file from the given path
	Exist(path string) bool
	// Return file content from the given path
	Retrieve(path string) (io.Reader, error)
	// Keep file content to the given path
	Store(body io.Reader, path string) error
}
