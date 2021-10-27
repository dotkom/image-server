package s3_adapter

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Implementation of storage.Storage using AWS S3 as its storage backend.
type S3Adapter struct {
	bucketName *string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	s3         *s3.S3
}

// Creates a new instance of S3Adapter taking in the bucket name and aws region as parameters
func New(bucketName string) (*S3Adapter, error) {
	adapter := &S3Adapter{}
	adapter.bucketName = aws.String(bucketName)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	})
	if err != nil {
		return nil, err
	}
	adapter.s3 = s3.New(sess)
	adapter.uploader = s3manager.NewUploader(sess)
	adapter.downloader = s3manager.NewDownloader(sess)
	return adapter, nil
}

func (adapter *S3Adapter) Save(ctx context.Context, content []byte, key string) error {
	reader := bytes.NewReader(content)
	_, err := adapter.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: adapter.bucketName,
		Key:    aws.String(key),
		Body:   reader,
	})
	return err
}

func (adapter *S3Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	buffer := &aws.WriteAtBuffer{}
	_, err := adapter.downloader.DownloadWithContext(ctx, buffer, &s3.GetObjectInput{
		Bucket: adapter.bucketName,
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (adapter *S3Adapter) Delete(ctx context.Context, key string) error {
	_, err := adapter.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: adapter.bucketName,
		Key:    aws.String(key),
	})

	return err
}
