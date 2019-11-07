package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"path/filepath"
)

// Service helps co-working between data-layer and control-layer
type Service interface {
	// A file repository
	File() FileRepository
}

type service struct {
	fileRepo FileRepository
}

func (s service) File() FileRepository {
	return s.fileRepo
}

// AmazonS3 is a storage manager which talks to Amazon S3 service for uploading and listing uploaded files
type AmazonS3 struct {
	service Service

	// Amazon Web Service session
	session *session.Session

	// The default bucket name requires when uploading file to Amazon S3
	defaultBucketName string
}

// NewAmazonS3 returns a new storage manager which configures default bucket name and session inside
func NewAmazonS3(accessKey, secretKey string, fileRepo FileRepository) (AmazonS3, error) {
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

	return AmazonS3{
		service: service{
			fileRepo: fileRepo,
		},
		session:           sess,
		defaultBucketName: "nomkhonwaan-com",
	}, nil
}

// SetDefaultBucketName allows to override the default bucket name
func (s AmazonS3) SetDefaultBucketName(newBucketName string) {
	s.defaultBucketName = newBucketName
}

func (s AmazonS3) Download(ctx context.Context, path string) (File, error) {
	file, err := s.service.File().FindByPath(ctx, path)
	if err != nil {
		return File{}, err
	}

	downloader := s3manager.NewDownloader(s.session)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.defaultBucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return File{}, err
	}

	file.Body = buf.Bytes()
	return file, nil
}

func (s AmazonS3) Upload(ctx context.Context, path string, body io.Reader) (File, error) {
	uploader := s3manager.NewUploader(s.session)

	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.defaultBucketName),
		Key:    aws.String(path),
		Body:   body,
	})
	if err != nil {
		return File{}, err
	}

	file, err := s.service.File().Create(ctx, File{
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
