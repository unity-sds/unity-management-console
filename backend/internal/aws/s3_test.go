package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/pkg/errors"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"strconv"
	"strings"
	"testing"
)

// type mockS3BucketAPI func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
type mockS3BucketAPI struct {
	CreateBucketFunc func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	HeadBucketFunc   func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
}

func (m mockS3BucketAPI) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.CreateBucketFunc(ctx, params, optFns...)
}

func (m mockS3BucketAPI) HeadBucket(ctx context.Context, params *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return m.HeadBucketFunc(ctx, params)
}
func TestCreateBucket(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) S3BucketAPI
		params s3.CreateBucketInput
		expect string
	}{
		{
			client: func(t *testing.T) S3BucketAPI {
				return mockS3BucketAPI{
					CreateBucketFunc: func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(options *s3.Options)) (*s3.CreateBucketOutput, error) {
						t.Helper()

						if params.Bucket == nil {
							t.Fatal("expect bucket to not be nil")
						}

						if e, a := "fooBucket", *params.Bucket; e != a {
							t.Errorf("expect %v, got %v", e, a)
						}

						return &s3.CreateBucketOutput{
							Location:       aws.String("fooBucket"),
							ResultMetadata: middleware.Metadata{},
						}, nil
					},
					HeadBucketFunc: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						if *params.Bucket == "fooBucket" {
							return &s3.HeadBucketOutput{}, nil
						}
						return nil, fmt.Errorf("bucket not found")
					},
				}
			},
			expect: "fooBucket",
			params: s3.CreateBucketInput{Bucket: aws.String("fooBucket")},
		},
	}

	for i, tt := range cases {
		//conf := config.AppConfig{AWSRegion: "us-west-2"}
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := CreateBucketFromS3(ctx, tt.client(t), &tt.params)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; strings.Compare(e, a) != 0 {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}

}

func TestHeadBucket(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) S3BucketAPI
		params s3.HeadBucketInput
		expect string
	}{
		{
			client: func(t *testing.T) S3BucketAPI {
				return mockS3BucketAPI{
					CreateBucketFunc: func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(options *s3.Options)) (*s3.CreateBucketOutput, error) {
						return &s3.CreateBucketOutput{
							Location:       aws.String("fooBucket"),
							ResultMetadata: middleware.Metadata{},
						}, nil
					},
					HeadBucketFunc: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						t.Helper()

						if *params.Bucket == "fooBucket" {
							return &s3.HeadBucketOutput{}, nil
						}
						return nil, fmt.Errorf("bucket not found")
					},
				}
			},
			expect: "fooBucket",
			params: s3.HeadBucketInput{Bucket: aws.String("fooBucket")},
		},
	}

	for i, tt := range cases {
		//conf := config.AppConfig{AWSRegion: "us-west-2"}
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := HeadBucketFromS3(ctx, tt.client(t), &tt.params)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}

}

func TestValidBucketCreation(t *testing.T) {
	api := &mockS3BucketAPI{
		CreateBucketFunc: func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
			t.Helper()
			var err error
			if params.Bucket == nil {
				t.Fatal("expect bucket to not be nil")
			}

			if !strings.HasPrefix(*params.Bucket, "mgmt") {
				t.Errorf("expect prefix mgmt, got %v", *params.Bucket)
				err = errors.New(fmt.Sprintf("expect prefix mgmt, got %v", *params.Bucket))
			}

			return &s3.CreateBucketOutput{
				Location:       aws.String("fooBucket"),
				ResultMetadata: middleware.Metadata{},
			}, err
		},
		HeadBucketFunc: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
			t.Helper()

			if *params.Bucket == "fooBucket" {
				return &s3.HeadBucketOutput{}, nil
			}
			return nil, fmt.Errorf("bucket not found")
		},
	}
	conf := config.AppConfig{AWSRegion: "us-west-2"}
	CreateBucket(api, &conf)
}

func TestBucketExistsBucketCreation(t *testing.T) {
	api := &mockS3BucketAPI{
		CreateBucketFunc: func(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
			t.Helper()
			var err error
			if params.Bucket == nil {
				t.Fatal("expect bucket to not be nil")
			}

			if !strings.HasPrefix(*params.Bucket, "mgmt") {
				t.Errorf("expect prefix mgmt, got %v", *params.Bucket)
				err = errors.New(fmt.Sprintf("expect prefix mgmt, got %v", *params.Bucket))
			}

			return &s3.CreateBucketOutput{
				Location:       aws.String("fooBucket"),
				ResultMetadata: middleware.Metadata{},
			}, err
		},
		HeadBucketFunc: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
			t.Helper()

			if strings.HasPrefix(*params.Bucket, "mgmt") {
				// Return no error as bucket already exists
				return &s3.HeadBucketOutput{}, nil
			}
			return nil, fmt.Errorf("bucket not found")
		},
	}
	conf := config.AppConfig{AWSRegion: "us-west-2"}
	CreateBucket(api, &conf)
}
