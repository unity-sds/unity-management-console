package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

func CreateBucket(conf *appconfig.AppConfig) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load AWS configuration: %v", err)
	}

	cfg.Region = conf.AWSRegion

	svc := s3.NewFromConfig(cfg)

	bucket := ""
	if conf.BucketName != "" {
		bucket = conf.BucketName
	} else {
		bucket = generateBucketName()
		conf.BucketName = bucket
		viper.Set("bucketname", bucket)
		err = viper.WriteConfigAs(viper.ConfigFileUsed())
		if err != nil {
			log.WithError(err).Error("Could not write config file")
		}
	}

	// Check if bucket exists
	_, err = svc.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		// Create bucket
		pass, perr := checkPolicy()
		if !pass || perr != nil {
			log.Warnf("Policy Check Failed, following actions may not work correctly. %v", err)
		}
		_, berr := svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
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
