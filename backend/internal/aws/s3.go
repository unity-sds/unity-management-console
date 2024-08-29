package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	appconfig "github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"io"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type S3BucketAPI interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error)
	PutBucketLifecycleConfiguration(ctx context.Context, params *s3.PutBucketLifecycleConfigurationInput) (*s3.PutBucketLifecycleConfigurationOutput, error)
}

type AWSS3Client struct {
	Client *s3.Client
}

func NewAWSS3Client(cfg aws.Config) S3BucketAPI {
	return &AWSS3Client{
		Client: s3.NewFromConfig(cfg),
	}
}

func (a *AWSS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return a.Client.CreateBucket(ctx, params, optFns...)
}

func (a *AWSS3Client) HeadBucket(ctx context.Context, params *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return a.Client.HeadBucket(ctx, params)
}

func (a *AWSS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return a.Client.GetObject(ctx, params)
}

func (a *AWSS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return a.Client.ListObjectsV2(ctx, params)
}

func (a *AWSS3Client) PutBucketVersioning(ctx context.Context, params *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return a.Client.PutBucketVersioning(ctx, params)
}

func (a *AWSS3Client) PutBucketLifecycleConfiguration(ctx context.Context, params *s3.PutBucketLifecycleConfigurationInput) (*s3.PutBucketLifecycleConfigurationOutput, error) {
	return a.Client.PutBucketLifecycleConfiguration(ctx, params)
}

func CreateBucketFromS3(ctx context.Context, api S3BucketAPI, params *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	resp, berr := api.CreateBucket(ctx, params)
	return resp, berr
}

func HeadBucketFromS3(ctx context.Context, api S3BucketAPI, params *s3.HeadBucketInput) error {
	_, err := api.HeadBucket(ctx, params)

	return err
}

func GetObjectFromS3(ctx context.Context, api S3BucketAPI, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return api.GetObject(ctx, params)
}

func ListObjectsFromS3(ctx context.Context, api S3BucketAPI, params *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(ctx, params)
}

func PutBucketVersioning(ctx context.Context, api S3BucketAPI, params *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	return api.PutBucketVersioning(ctx, params)
}

func PutBucketLifecycleConfiguration(ctx context.Context, api S3BucketAPI, params *s3.PutBucketLifecycleConfigurationInput) (*s3.PutBucketLifecycleConfigurationOutput, error) {
	return api.PutBucketLifecycleConfiguration(ctx, params)
}

func InitS3Client(conf *appconfig.AppConfig) S3BucketAPI {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(conf.AWSRegion))
	if err != nil {
		log.Fatalf("failed to load AWS configuration: %v", err)
	}
	return NewAWSS3Client(cfg)
}

func CreateBucket(s3client S3BucketAPI, conf *appconfig.AppConfig) {
	if s3client == nil {
		s3client = InitS3Client(conf)
	}

	bucket := ""
	if conf.BucketName != "" {
		bucket = conf.BucketName
	} else {
		// We get the name of this deployment's bucket from an SSM param (previously was a randomized string prefix)
		bucketNameParamPath := fmt.Sprintf("/unity/%s/%s/cs/monitoring/s3/bucketName", conf.Project, conf.Venue)
		bucketNameParam, err := ReadSSMParameter(bucketNameParamPath)

		if err != nil {
			log.WithError(err).Error("Could not find SSM parameter for bucket name at: %s", bucketNameParamPath)
		}

		bucket = *bucketNameParam.Parameter.Value
		conf.BucketName = bucket
		viper.Set("bucketname", bucket)
		err = viper.WriteConfigAs(viper.ConfigFileUsed())
		if err != nil {
			log.WithError(err).Error("Could not write config file")
		}
	}

	// Check if bucket exists
	bucketinput := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}
	err := HeadBucketFromS3(context.TODO(), s3client, bucketinput)

	if err != nil {
		// Create bucket
		pass, perr := checkPolicy()
		if !pass || perr != nil {
			log.Warnf("Policy Check Failed, following actions may not work correctly. %v", err)
		}

		_, berr := CreateBucketFromS3(context.TODO(), s3client, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(conf.AWSRegion),
			},
		})

		if berr != nil {
			log.Errorf("Error creating bucket: %v", berr)
			return
		}

		// Enable versioning on bucket
		berr = EnableBucketVersioning(s3client, conf, bucket)

		// Set bucket lifecycle length
		berr = SetBucketLifecycleLength(s3client, conf, bucket, int32(7))

		if berr != nil {
			log.Errorf("Error enabling versioning on bucket: %v", berr)
			return
		}
	} else {
		log.Infof("Bucket %s exists", bucket)
	}
}

// Empties, then deletes a S3 bucket
func DeleteS3Bucket(bucketName string) error {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("us-west-2"))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	bucket := aws.String(bucketName)
	client := s3.NewFromConfig(cfg)

	// Define an inline deleteObject function (used below)
	deleteObject := func(bucket, key, versionId *string) {
		log.Printf("Object: %s/%s\n", *key, aws.ToString(versionId))
		_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: versionId,
		})
		if err != nil {
			log.Fatalf("Failed to delete object: %v", err)
		}
	}

	//
	// Get a list of objects in the S3 bucket.
	// Iterator over them, and delete each one.
	//
	in := &s3.ListObjectsV2Input{Bucket: bucket}
	for {
		out, err := client.ListObjectsV2(context.TODO(), in)
		if err != nil {
			log.Fatalf("Failed to list objects: %v", err)
		}

		for _, item := range out.Contents {
			deleteObject(bucket, item.Key, nil)
		}

		if out.IsTruncated {
			in.ContinuationToken = out.ContinuationToken
		} else {
			break
		}
	}

	//
	// Get a list of all the object versions in the bucket.
	// Iterate over them and delete them
	//
	inVer := &s3.ListObjectVersionsInput{Bucket: bucket}
	for {
		out, err := client.ListObjectVersions(context.TODO(), inVer)
		if err != nil {
			log.Fatalf("Failed to list version objects: %v", err)
		}

		for _, item := range out.DeleteMarkers {
			deleteObject(bucket, item.Key, item.VersionId)
		}

		for _, item := range out.Versions {
			deleteObject(bucket, item.Key, item.VersionId)
		}

		if out.IsTruncated {
			inVer.VersionIdMarker = out.NextVersionIdMarker
			inVer.KeyMarker = out.NextKeyMarker
		} else {
			break
		}
	}

	//
	// Finally, delete the (now empty) bucket
	//
	_, err = client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetObject(s3client S3BucketAPI, conf *appconfig.AppConfig, bucketName string, objectKey string) ([]byte, error) {
	if s3client == nil {
		s3client = InitS3Client(conf)
	}

	objectinput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	result, err := GetObjectFromS3(context.TODO(), s3client, objectinput)

	// defer result.Body.Close()

	if err != nil {
		log.WithError(err).Error("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}

	bytesRead, err := io.ReadAll(result.Body)
	if err != nil {
		log.WithError(err).Error("Unable to read object")
		return nil, err
	}

	return bytesRead, nil
}

func ListObjectsV2(s3client S3BucketAPI, conf *appconfig.AppConfig, bucketName string, prefix string) []types.Object {
	if s3client == nil {
		s3client = InitS3Client(conf)
	}

	listobjectsinput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	result, err := ListObjectsFromS3(context.TODO(), s3client, listobjectsinput)

	if err != nil {
		log.WithError(err).Error("Couldn't list objects in bucket: %s. Here's why: %v\n", bucketName, err)
	}

	return result.Contents

}

func EnableBucketVersioning(s3client S3BucketAPI, conf *appconfig.AppConfig, bucketName string) error {
	if s3client == nil {
		s3client = InitS3Client(conf)
	}

	putBucketVersioningInput := &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: "Enabled",
		},
	}

	_, err := PutBucketVersioning(context.TODO(), s3client, putBucketVersioningInput)

	if err != nil {
		log.WithError(err).Error("Couldn't enable bucket versioning: %s. Here's why: %v\n", bucketName, err)
	}

	return nil
}

func SetBucketLifecycleLength(s3client S3BucketAPI, conf *appconfig.AppConfig, bucketName string, lifecycleInDays int32) error {
	if s3client == nil {
		s3client = InitS3Client(conf)
	}

	lifecycleRule := &types.LifecycleRule{
		Expiration: &types.LifecycleExpiration{
			Days: lifecycleInDays,
		},
		Prefix: nil,
	}

	putBucketLifecycleConfigurationInput := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{
				*lifecycleRule,
			},
		},
	}

	_, err := PutBucketLifecycleConfiguration(context.TODO(), s3client, putBucketLifecycleConfigurationInput)

	if err != nil {
		log.WithError(err).Error("Couldn't enable bucket lifecycle configuration for: %s. Here's why: %v\n", bucketName, err)
	}

	return nil
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func stringf(length int) string {
	return stringWithCharset(length, charset)
}

func generateBucketName() string {
	return "mgmt-" + stringf(8)
}
