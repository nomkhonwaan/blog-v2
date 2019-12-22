package gcloud

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
)

// CloudStorage implements Uploader, Downloader and Cache from storage package with Google Cloud Storage service
type CloudStorage struct {
	credentialsFile string
	bucket          string
}

// NewCloudStorage returns new CloudStorage instance
func NewCloudStorage(credentialsFilePath, bucket string) CloudStorage {
	return CloudStorage{credentialsFile: credentialsFilePath, bucket: bucket}
}

func (cs CloudStorage) newClient(ctx context.Context) (*storage.Client, error) {
	opts := make([]option.ClientOption, 0)

	if cs.credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cs.credentialsFile))
	}

	return storage.NewClient(ctx, opts...)
}
func (cs CloudStorage) Download(ctx context.Context, path string) (io.Reader, error) {
	client, err := cs.newClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.Bucket(cs.bucket).Object(path).NewReader(ctx)
}

func (cs CloudStorage) Upload(ctx context.Context, body io.Reader, path string) error {
	client, err := cs.newClient(ctx)
	if err != nil {
		return err
	}

	wtr := client.Bucket(cs.bucket).Object(path).NewWriter(ctx)
	if _, err = io.Copy(wtr, body); err != nil {
		return err
	}

	return wtr.Close()
}

//func (cs CloudStorage) Exist(path string) bool {
//	client, err := cs.newClient(context.Background())
//	if err != nil {
//		return false
//	}
//
//	_, err = client.Bucket(cs.bucket).Object(path).Attrs(context.Background())
//	if err != nil {
//		return false
//	}
//
//	return true
//}
//
//func (cs CloudStorage) Retrieve(path string) (io.Reader, error) {
//	return cs.Download(context.Background(), path)
//}
//
//func (cs CloudStorage) Store(body io.Reader, path string) error {
//	return cs.Upload(context.Background(), body, path)
//}
