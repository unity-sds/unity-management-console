package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	appconfig "github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type S3BucketAPI interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
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

func CreateBucketFromS3(ctx context.Context, api S3BucketAPI, params *s3.CreateBucketInput) (string, error) {
	resp, berr := api.CreateBucket(ctx, params)

	return *resp.Location, berr
}

func HeadBucketFromS3(ctx context.Context, api S3BucketAPI, params *s3.HeadBucketInput) error {
	_, err := api.HeadBucket(ctx, params)

	return err
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
		bucket = generateBucketName()
		conf.BucketName = bucket
		viper.Set("bucketname", bucket)
		err := viper.WriteConfigAs(viper.ConfigFileUsed())
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
	} else {
		log.Infof("Bucket %s exists", bucket)
	}
}

func DeleteS3Bucket(bucketName string) error {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("us-west-2"))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
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
