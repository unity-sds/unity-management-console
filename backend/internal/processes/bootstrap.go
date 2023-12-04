package processes

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"path/filepath"
)

func BootstrapEnv(appconf *config.AppConfig) {
	store, err := database.NewGormDatastore()

	provisionS3(appconf)
	initTerraform(store, appconf)

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

func initTerraform(store database.Datastore, appconf *config.AppConfig) {
	fs := afero.NewOsFs()
	writeInitTemplate(fs, appconf)
	installUnityCloudEnv(store, appconf)

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

	err := fs.MkdirAll(filepath.Join(appConfig.Workdir, "workspace"), 0755)
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
		Name:        "unity-proxy",
		Version:     "0.1",
		Variables:   nil,
		Displayname: "unity-proxy",
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

func installUnityCloudEnv(store database.Datastore, appConfig *config.AppConfig) error {

	venue, err := getSSMParameterValueFromDatabase("venue", store)
	if err != nil {
		log.WithError(err).Error("Problem fetching venue")
		return err
	}
	log.Infof("Venue found: %s", venue)
	project, err := getSSMParameterValueFromDatabase("project", store)
	if err != nil {
		log.WithError(err).Error("Problem fetching project")
		return err
	}
	log.Infof("Project found: %s", project)

	publicsubnets, err := getSSMParameterValueFromDatabase("publicsubnets", store)
	if err != nil {
		log.WithError(err).Error("Problem fetching public subnets")
		return err
	}
	log.Infof("Public subnets found: %s", publicsubnets)
	privatesubnets, err := getSSMParameterValueFromDatabase("privatesubnets", store)
	if err != nil {
		log.WithError(err).Error("Problem fetching private subnets")
		return err
	}
	log.Infof("Private subnets found: %s", privatesubnets)

	//ssmParameters, err := generateSSMParameters(db)
	//if err != nil {
	//	log.WithError(err).Error("Problem fetching params")
	//	return err
	//}

	varmap := make(map[string]string)

	varmap["venue"] = venue
	varmap["project"] = project
	varmap["publicsubnets"] = publicsubnets
	varmap["privatesubnets"] = privatesubnets
	vars := marketplace.Install_Variables{
		Values:         varmap,
		AdvancedValues: nil,
	}
	applications := marketplace.Install_Applications{
		Name:        "unity-cloud-env",
		Version:     "0.1",
		Variables:   &vars,
		Displayname: "unity-cloud-env",
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Unity Cloud Environment",
	}
	err = TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing Unity Cloud Env")
		return err
	}
	return nil
}
