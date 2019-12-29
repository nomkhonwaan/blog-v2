package aws

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

// S3 implements Uploader, Downloader from storage package with Amazon S3 service
type S3 struct {
	*session.Session

	bucket string
}

// NewS3 returns new S3 instance
func NewS3(accessKey, secretKey, bucket string) (S3, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Region:      aws.String("ap-southeast-1"),
		},
	})
	if err != nil {
		return S3{}, err
	}

	return S3{
		Session: s,
		bucket:  bucket,
	}, nil
}

func (s S3) Delete(ctx context.Context, path string) error {
	svc := s3.New(s.Session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}

	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
}

func (s S3) Download(ctx context.Context, path string) (io.Reader, error) {
	downloader := s3manager.NewDownloader(s.Session)
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func (s S3) Upload(ctx context.Context, body io.Reader, path string) error {
	uploader := s3manager.NewUploader(s.Session)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   body,
	})

	return err
}
