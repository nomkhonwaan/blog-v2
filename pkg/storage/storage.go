//go:generate mockgen -destination=./mock/storage_mock.go github.com/nomkhonwaan/myblog/pkg/storage Storage

package storage

import (
	"context"
	"io"
)

// Storage uses to storing or retrieving file from cloud or remote server
type Storage interface {
	Delete(ctx context.Context, path string) error
	Download(ctx context.Context, path string) (io.Reader, error)
	Upload(ctx context.Context, body io.Reader, path string) error
}
