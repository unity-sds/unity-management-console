package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockS3Client struct {
	s3.Client

	GetObjectOutput *s3.GetObjectOutput
}

func (m *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectOutput, nil
}

type mockS3Client struct {
	s3.Client
	createBucketFunc func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

func (m *mockS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.createBucketFunc(ctx, params, optFns...)
}

func TestCreateBucket(t *testing.T) {

}
