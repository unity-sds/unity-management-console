package action

import (
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"testing"
)

func TestRunInstall(t *testing.T) {
	nodeGroup := &marketplace.Install_Extensions_Nodegroups{
		Name:         "default",
		Instancetype: "m6.xlarge",
		Nodecount:    "10",
	}
	eks := &marketplace.Install_Extensions_Eks{
		Clustername: "test cluster",
		Owner:       "tom",
		Projectname: "test",
		Nodegroups:  []*marketplace.Install_Extensions_Nodegroups{nodeGroup},
	}
	extensions := &marketplace.Install_Extensions{
		Eks:        eks,
		Apigateway: nil,
	}
	install := &marketplace.Install{
		Applications:   nil,
		Extensions:     extensions,
		DeploymentName: "my deployment",
	}

	appConfig := config.AppConfig{
		GithubToken:          "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",
		WorkflowBasePath:     "/home/barber/Projects/unity-cs-infra/.github/workflows",
		AWSRegion:            "",
		BucketName:           "",
		Workdir:              "",
		DefaultSSMParameters: nil,
	}
	r := ActRunnerImpl{}
	err := RunInstall(install, nil, appConfig, r)
	if err != nil {
		return
	}
}
