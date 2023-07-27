package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	log "github.com/sirupsen/logrus"
)
import "github.com/aws/aws-sdk-go-v2/service/iam"

func checkPolicy() (bool, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))

	if err != nil {
		log.Errorf("Error creating session: %v", err)
		return false, err
	}

	client := iam.NewFromConfig(cfg)

	resp, err := client.SimulatePrincipalPolicy(context.TODO(), &iam.SimulatePrincipalPolicyInput{
		ActionNames:                        nil,
		PolicySourceArn:                    nil,
		CallerArn:                          nil,
		ContextEntries:                     nil,
		Marker:                             nil,
		MaxItems:                           nil,
		PermissionsBoundaryPolicyInputList: nil,
		PolicyInputList:                    nil,
		ResourceArns:                       nil,
		ResourceHandlingOption:             nil,
		ResourceOwner:                      nil,
		ResourcePolicy:                     nil,
	})

	log.Infof("IAM Response: %v", resp)

	//TODO return err
	return true, nil
}
