package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
)

func ReadSSMParameter(path string) (*ssm.GetParameterOutput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ssm.NewFromConfig(cfg)

	input := &ssm.GetParameterInput{
		Name: &path,
	}

	return client.GetParameter(context.TODO(), input)
}
func ReadSSMParameters(ssmParams []models.SSMParameters) (*marketplace.Parameters, error) {
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
	ssmparams := map[string]*marketplace.Parameters_Parameter{}
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
			exists, valuetracks := checkExistsInDatabase(aws.ToString(param.Name), aws.ToString(param.Value), ssmParams)
			par := marketplace.Parameters_Parameter{
				Value:   aws.ToString(param.Value),
				Type:    "fixme",
				Tracked: exists,
				Insync:  valuetracks,
			}
			ssmparams[aws.ToString(param.Name)] = &par
		}
	}

	ssmparamholder.Parameterlist = ssmparams
	return &ssmparamholder, nil
}

func checkExistsInDatabase(name, value string, params []models.SSMParameters) (bool, bool) {

	for _, param := range params {
		if param.Key == name {
			if param.Value == value {
				return true, true
			} else {
				return true, false
			}
		}
	}

	return false, false
}
