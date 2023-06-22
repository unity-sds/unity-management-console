package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
)

func ReadSSMParameters() (*marketplace.Parameters, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ssm.NewFromConfig(cfg)

	recurse := true
	decrypt := true
	paramsInput := &ssm.GetParametersByPathInput{
		Path:           aws.String("/unity/"), // replace with the actual path
		Recursive:      &recurse,
		WithDecryption: &decrypt,
	}

	paginator := ssm.NewGetParametersByPathPaginator(client, paramsInput)
	ssmparams := map[string]string{}
	ssmparamholder := marketplace.Parameters{}

	// Iterate through the SSM parameter pages.
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.WithError(err).Error("failed to get parameters")
			return &ssmparamholder, err
		}

		// Print out the parameters.
		for _, param := range output.Parameters {
			log.Infof("Name: %s, Value: %s\n", aws.ToString(param.Name), aws.ToString(param.Value))
			ssmparams[aws.ToString(param.Name)] = aws.ToString(param.Value)
		}
	}

	ssmparamholder.Parameterlist = ssmparams
	return &ssmparamholder, nil
}
