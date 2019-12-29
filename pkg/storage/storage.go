// go:generate mockgen -destination=./mock/storage_mock.go github.com/nomkhonwaan/myblog/storage Storage

package storage

import (
	"context"
	"io"
)

// Storage keeps data at the storage service where application can download, upload and delete from everywhere
type Storage interface {
	// Permanently remove uploaded file from the storage server
	Delete(ctx context.Context, path string) error
	// Retrieve uploaded file from the storage server
	Download(ctx context.Context, path string) (io.Reader, error)
	// Push file to the storage server
	Upload(ctx context.Context, body io.Reader, path string) error
}
