package processes

import (
	"fmt"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	config2 "github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"testing"
)

func TestInstallMarketplaceApplication(t *testing.T) {

	appConf := config2.AppConfig{
		GithubToken:          "",
		MarketplaceOwner:     "",
		MarketplaceRepo:      "",
		AWSRegion:            "",
		BucketName:           "mgmt-jysn71x6",
		Workdir:              "/home/barber/Projects/unity-management-console/workdir",
		DefaultSSMParameters: nil,
	}

	apps := marketplace.Install_Applications{
		Name:      "unity-apigateway",
		Version:   "0.1",
		Variables: nil,
	}
	install := marketplace.Install{
		Applications:   &apps,
		DeploymentName: "test",
	}

	store, _ := database.NewGormDatastore()

	err := TriggerInstall(nil, "", store, &install, &appConf)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}

func TestUpdateCoreConfig(t *testing.T) {
	store, _ := database.NewGormDatastore()
	appConf := config2.AppConfig{
		GithubToken:          "",
		MarketplaceOwner:     "",
		MarketplaceRepo:      "",
		AWSRegion:            "",
		BucketName:           "mgmt-jysn71x6",
		Workdir:              "/home/barber/Projects/unity-management-console/workdir",
		DefaultSSMParameters: nil,
	}

	p := []config2.SSMParameter{{
		Name:  "venue",
		Type:  "String",
		Value: "test",
	}, {
		Name:  "project",
		Type:  "String",
		Value: "unity-nightly",
	}}
	store.StoreSSMParams(p, "test")
	err := UpdateCoreConfig(&appConf, store, nil, "")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
