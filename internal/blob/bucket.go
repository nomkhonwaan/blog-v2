package blob

import (
	"context"
	"gocloud.dev/blob"
	"io"
)

// Bucket embeds the original blob.Bucket for providing compatible storage.Storage interface methods
type Bucket struct{ *blob.Bucket }

// Download provides compatible Download method of the storage.Storage interface
func (b *Bucket) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	return b.Bucket.NewReader(ctx, path, nil)
}

// Upload provides compatible Upload method of the storage.Storage interface
func (b *Bucket) Upload(ctx context.Context, body io.Reader, path string) error {
	w, err := b.Bucket.NewWriter(ctx, path, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, body)
	return err
}
