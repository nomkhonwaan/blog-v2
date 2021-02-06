// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package app

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nomkhonwaan/myblog/internal/blob"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	blob2 "gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

import (
	_ "gocloud.dev/blob/s3blob"
)

// Injectors from provide_storage.go:

func provideS3Storage() (storage.Storage, error) {
	bucket, err := provideS3GoCloudBlobBucket()
	if err != nil {
		return nil, err
	}
	blobBucket := &blob.Bucket{
		Bucket: bucket,
	}
	return blobBucket, nil
}

func provideS3GoCloudBlobBucket() (*blob2.Bucket, error) {
	contextContext := context.Background()
	session, err := provideS3Session()
	if err != nil {
		return nil, err
	}
	string2 := provideS3BucketName()
	options := provideS3BucketOptions()
	bucket, err := s3blob.OpenBucket(contextContext, session, string2, options)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func provideS3Session() (*session.Session, error) {
	options := provideS3SessionOptions()
	sessionSession, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, err
	}
	return sessionSession, nil
}

func provideS3SessionOptions() session.Options {
	config := provideS3Config()
	options := session.Options{
		Config: config,
	}
	return options
}

func provideS3Config() aws.Config {
	credentials := provideS3StaticCredentials()
	string2 := provideS3Region()
	config := aws.Config{
		Credentials: credentials,
		Region:      string2,
	}
	return config
}

func provideDiskCache(fs afero.Fs) (storage.Cache, error) {
	string2 := provideDiskCacheFilePath()
	diskCache, err := storage.NewDiskCache(fs, string2)
	if err != nil {
		return nil, err
	}
	return diskCache, nil
}

// provide_storage.go:

func provideS3BucketName() string {
	return viper.GetString("amazon-s3-bucket-name")
}

func provideS3BucketOptions() *s3blob.Options {
	return nil
}

func provideS3StaticCredentials() *credentials.Credentials {
	return credentials.NewStaticCredentials(viper.GetString("amazon-s3-access-key"), viper.GetString("amazon-s3-secret-key"), "")
}

func provideS3Region() *string {
	return aws.String(viper.GetString("amazon-s3-region"))
}

func provideDiskCacheFilePath() string {
	return viper.GetString("cache-file-path")
}
