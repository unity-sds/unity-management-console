package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func listObjects(client *s3.Client, bucket string) (*s3.ListObjectsV2Output, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: 7,
	}
	return client.ListObjectsV2(context.TODO(), input)
}
