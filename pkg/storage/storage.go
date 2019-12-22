package storage

import (
	"context"
	"io"
)

// Downloader downloads file from the storage server
type Downloader interface {
	Download(ctx context.Context, path string) (io.Reader, error)
}

// Uploader uploads file to the storage server
type Uploader interface {
	Upload(ctx context.Context, body io.Reader, path string) error
}
