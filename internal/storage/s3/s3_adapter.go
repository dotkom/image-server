package s3_adapter

import (
	"bytes"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Implementation of storage.Storage using AWS S3 as its storage backend.
type S3Adapter struct {
	bucketName string
	client     *s3.Client
}

// Creates a new instance of S3Adapter taking in the bucket name and aws region as parameters
func New(bucketName string) (*S3Adapter, error) {
	adapter := &S3Adapter{}

	if bucketName == "" {
		return nil, errors.New("s3 bucket name cannot be empty")
	}
	adapter.bucketName = bucketName

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	adapter.client = s3.NewFromConfig(cfg)
	return adapter, nil
}

func (adapter *S3Adapter) Save(ctx context.Context, content []byte, key string) error {
	uploader := manager.NewUploader(adapter.client)
	reader := bytes.NewReader(content)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &adapter.bucketName,
		Key:    &key,
		Body:   reader,
	})
	return err
}

func (adapter *S3Adapter) Get(ctx context.Context, key string, size uint64) ([]byte, error) {
	downloader := manager.NewDownloader(adapter.client)
	buffer := make([]byte, int(size))
	writer := manager.NewWriteAtBuffer(buffer)
	_, err := downloader.Download(ctx, writer, &s3.GetObjectInput{
		Bucket: &adapter.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func (adapter *S3Adapter) Delete(ctx context.Context, key string) error {
	_, err := adapter.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &adapter.bucketName,
		Key:    &key,
	})

	return err
}
