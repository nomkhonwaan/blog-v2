package storage

import (
	"context"
	"io"
)

// Downloader uses to downloading file from the storage server
type Downloader interface {
	Download(ctx context.Context, path string) (File, error)
}

// Uploader uses to uploading file from multipart body to the storage server
type Uploader interface {
	Upload(ctx context.Context, path string, body io.Reader) (File, error)
}
