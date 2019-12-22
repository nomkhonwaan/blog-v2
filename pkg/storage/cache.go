package storage

import (
	"io"
)

// Cache service will be implemented this interface for manipulating cache file with its storage
type Cache interface {
	Exist(path string) bool
	Retrieve(path string) (io.Reader, error)
	Store(body io.Reader, path string) error
}
