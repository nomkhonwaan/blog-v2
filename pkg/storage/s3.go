package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"path/filepath"
)

// AmazonS3 is a storage manager which talks to Amazon S3 service for uploading and listing uploaded files
type AmazonS3 struct {
	repo FileRepository
	sess *session.Session

	// The default bucket name requires when uploading file to Amazon S3 which keep all files in the bucket storage
	defaultBucketName string
}

// NewAmazonS3 returns a new storage manager which configures default bucket name and session inside
func NewAmazonS3(accessKey, secretKey string, repo FileRepository) (AmazonS3, error) {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
				Region:      aws.String("ap-southeast-1"),
			},
		},
	)
	if err != nil {
		return AmazonS3{}, err
	}

	return AmazonS3{repo: repo, sess: sess, defaultBucketName: "nomkhonwaan-com"}, nil
}

// SetDefaultBucketName allows to override the default bucket name
func (s AmazonS3) SetDefaultBucketName(newBucketName string) {
	s.defaultBucketName = newBucketName
}

func (s AmazonS3) Upload(ctx context.Context, path string, body io.Reader) (File, error) {
	u := s3manager.NewUploader(s.sess)

	result, err := u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.defaultBucketName),
		Key:    aws.String(path),
		Body:   body,
	})
	if err != nil {
		return File{}, err
	}

	file, err := s.repo.Create(ctx, File{
		Path:           path,
		FileName:       filepath.Base(path),
		OptionalField1: fmt.Sprintf("%T", s),
		OptionalField2: result.UploadID,
		OptionalField3: result.Location,
	})
	if err != nil {
		return File{}, err
	}

	return file, nil
}
