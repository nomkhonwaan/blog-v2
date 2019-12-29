package gcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
)

// CloudStorage implements Uploader, Downloader and Cache from storage package with Google Cloud Storage service
type CloudStorage struct {
	*storage.BucketHandle
}

// NewCloudStorage returns new CloudStorage instance
func NewCloudStorage(credentialsFilePath, bucket string) (CloudStorage, error) {
	opts := make([]option.ClientOption, 0)

	if credentialsFilePath != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFilePath))
	}

	client, err := storage.NewClient(context.Background(), opts...)
	if err != nil {
		return CloudStorage{}, err
	}

	return CloudStorage{client.Bucket(bucket)}, nil
}

func (cs CloudStorage) Delete(ctx context.Context, path string) error {
	return cs.BucketHandle.Object(path).Delete(ctx)
}

func (cs CloudStorage) Download(ctx context.Context, path string) (io.Reader, error) {
	return cs.BucketHandle.Object(path).NewReader(ctx)
}

func (cs CloudStorage) Upload(ctx context.Context, body io.Reader, path string) error {
	wtr := cs.BucketHandle.Object(path).NewWriter(ctx)
	if _, err := io.Copy(wtr, body); err != nil {
		return err
	}

	return wtr.Close()
}
