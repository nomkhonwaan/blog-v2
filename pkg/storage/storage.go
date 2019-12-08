package storage

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

// Downloader uses to downloading file from the storage server
type Downloader interface {
	Download(ctx context.Context, path string) (io.Reader, error)
}

// Uploader uses to uploading file from multipart body to the storage server
type Uploader interface {
	Upload(ctx context.Context, body io.Reader, path string) error
}

// CustomizedAmazonS3Client is an implementation of Downloader and Uploader interfaces for using Amazon S3 as a storage
type CustomizedAmazonS3Client struct {
	session *session.Session

	// A bucket is required when downloading or uploading file to the Amazon S3 service
	defaultBucketName string
}

// NewCustomizedAmazonS3Client returns a new Amazon S3 client instance
func NewCustomizedAmazonS3Client(accessKey, secretKey string) (CustomizedAmazonS3Client, error) {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
				Region:      aws.String("ap-southeast-1"),
			},
		},
	)
	if err != nil {
		return CustomizedAmazonS3Client{}, err
	}

	return CustomizedAmazonS3Client{
		session:           sess,
		defaultBucketName: "nomkhonwaan-com",
	}, nil
}

// SetDefaultBucketName allows to override the default bucket name
func (s CustomizedAmazonS3Client) SetDefaultBucketName(newBucketName string) {
	s.defaultBucketName = newBucketName
}

func (s CustomizedAmazonS3Client) Download(ctx context.Context, path string) (io.Reader, error) {
	downloader := s3manager.NewDownloader(s.session)
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.defaultBucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func (s CustomizedAmazonS3Client) Upload(ctx context.Context, body io.Reader, path string) error {
	uploader := s3manager.NewUploader(s.session)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.defaultBucketName),
		Key:    aws.String(path),
		Body:   body,
	})
	return err
}
