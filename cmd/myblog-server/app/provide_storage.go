//+build wireinject

package app

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/wire"
	"github.com/nomkhonwaan/myblog/internal/blob"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	gocloud_blob "gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

func provideS3Storage() (storage.Storage, error) {
	panic(wire.Build(
		provideS3GoCloudBlobBucket,
		wire.Struct(new(blob.Bucket), "Bucket"),
		wire.Bind(new(storage.Storage), new(*blob.Bucket)),
	))
	return nil, nil
}

func provideS3GoCloudBlobBucket() (*gocloud_blob.Bucket, error) {
	panic(wire.Build(
		context.Background,
		provideS3Session,
		provideS3BucketName,
		provideS3BucketOptions,
		s3blob.OpenBucket,
		wire.Bind(new(client.ConfigProvider), new(*session.Session)),
	))
	return nil, nil
}

func provideS3BucketName() string {
	return viper.GetString("amazon-s3-bucket-name")
}

func provideS3BucketOptions() *s3blob.Options {
	return nil
}

func provideS3Session() (*session.Session, error) {
	panic(wire.Build(
		provideS3SessionOptions,
		session.NewSessionWithOptions,
	))
	return nil, nil
}

func provideS3SessionOptions() session.Options {
	panic(wire.Build(
		provideS3Config,
		wire.Struct(new(session.Options), "Config"),
	))
	return session.Options{}
}

func provideS3Config() aws.Config {
	panic(wire.Build(
		provideS3StaticCredentials,
		provideS3Region,
		wire.Struct(new(aws.Config), "Credentials", "Region"),
	))
	return aws.Config{}
}

func provideS3StaticCredentials() *credentials.Credentials {
	return credentials.NewStaticCredentials(viper.GetString("amazon-s3-access-key"),
		viper.GetString("amazon-s3-secret-key"), "")
}

func provideS3Region() *string {
	return aws.String(viper.GetString("amazon-s3-region"))
}

func provideDiskCache(fs afero.Fs) (storage.Cache, error) {
	panic(wire.Build(
		provideDiskCacheFilePath,
		storage.NewDiskCache,
		wire.Bind(new(storage.Cache), new(*storage.DiskCache)),
	))

	return nil, nil
}

func provideDiskCacheFilePath() string {
	return viper.GetString("cache-file-path")
}
