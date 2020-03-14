//go:generate mockgen -destination=./mock/cache_mock.go github.com/nomkhonwaan/myblog/pkg/storage Cache

package storage

import (
	"io"
)

// Cache uses to storing or retrieving files from hidden or inaccessible place
type Cache interface {
	Exist(path string) bool
	Retrieve(path string) (io.Reader, error)
	Store(body io.Reader, path string) error
}
