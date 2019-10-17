package storage

import (
	"context"
	"io"
)

// Uploader is the interface that wraps the file uploading method.
type Uploader interface {
	// Perform file uploading to the destination whether cloud service or internal storage
	Upload(ctx context.Context, path string, body io.Reader) (File, error)
}
