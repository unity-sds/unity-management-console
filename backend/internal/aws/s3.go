package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func CreateBucket(conf config.AppConfig){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(conf.AWSRegion),
	})

	if err != nil {
		log.Errorf("Error creating session: %v", err)
		return
	}

	svc := s3.New(sess)

	bucket := ""
	if conf.BucketName != "" {
		bucket = conf.BucketName
	} else {
		bucket = generateBucketName()
		//TODO Persist bucketname in config
	}

	// check bucket exists
	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		// create bucket
		pass, err := checkPolicy()
		if pass == false || err != nil {
			log.Warnf("Policy Check Failed, following actions may not work correctly. %v", err)
		}
		_, err = svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})

		if err != nil {
			log.Errorf("Error creating bucket, %v", err)
			return
		}
	} else {
		log.Infof("Bucket %v exists", bucket)
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
	return "mgmt_" + stringf(8)
}