package processes

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"path/filepath"
)

func BootstrapEnv(appconf *config.AppConfig) {
	store, err := database.NewGormDatastore()

	provisionS3(appconf)
	initTerraform(appconf)

	storeDefaultSSMParameters(appconf, store)
	//r := action.ActRunnerImpl{}
	err = UpdateCoreConfig(appconf, store, nil, "")
	if err != nil {
		log.WithError(err).Error("Problem updating ssm config")
	}
	installGateway(store, appconf)
}

func provisionS3(appConfig *config.AppConfig) {
	aws.CreateBucket(nil, appConfig)
	err := aws.CreateTable(appConfig)
	if err != nil {
		log.WithError(err).Error("Error creating table")
		return
	}
}

func initTerraform(appconf *config.AppConfig) {
	executor := &terraform.RealTerraformExecutor{}
	fs := afero.NewOsFs()
	writeInitTemplate(fs, appconf)
	err := terraform.RunTerraform(appconf, nil, "", executor, "unity-cloud-env")
	if err != nil {
		return
	}

}

func writeInitTemplate(fs afero.Fs, appConfig *config.AppConfig) {
	// Define the terraform configuration
	tfconfig := `terraform {
required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    dynamodb_table = "terraform_state"
  }
}

provider "aws" {
  region = "us-west-2"
}`

	err := fs.MkdirAll(appConfig.Workdir, 0755)
	if err != nil {
		log.WithError(err).Error("Couldn't create new working directory")
	}

	// Create a new file
	file, err := fs.Create(filepath.Join(appConfig.Workdir, "workspace", "provider.tf"))
	if err != nil {
		log.WithError(err).Error("Couldn't create new provider.tf file")
	}
	defer file.Close()

	// Write the configuration to the file
	_, err = file.WriteString(tfconfig)
	if err != nil {
		log.WithError(err).Error("Could not write provider string")
	}

	// Save changes to the file
	err = file.Sync()
	if err != nil {
		log.WithError(err).Error("Could not write provider.tf")
	}

	log.Println("File 'provider.tf' has been written")
}
func storeDefaultSSMParameters(appConfig *config.AppConfig, store database.Datastore) {

	err := store.StoreSSMParams(appConfig.DefaultSSMParameters, "bootstrap")
	if err != nil {
		log.WithError(err).Error("Problem storing parameters in database.")
	}
}

func installGateway(store database.Datastore, appConfig *config.AppConfig) {
	applications := marketplace.Install_Applications{
		Name:      "unity-proxy",
		Version:   "0.1",
		Variables: nil,
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Core API Gateway",
	}
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing API Gateway")
	}
}
