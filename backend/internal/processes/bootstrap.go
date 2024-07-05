package processes

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"path/filepath"
	"strings"
)

func BootstrapEnv(appconf *config.AppConfig) {
	store, err := database.NewGormDatastore()
	if err != nil {
		log.WithError(err).Error("Problem creating database")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	err = provisionS3(appconf)
	if err != nil {
		log.WithError(err).Error("Error provisioning S3 bucket")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}
	err = storeDefaultSSMParameters(appconf, store)
	if err != nil {
		log.WithError(err).Error("Error setting SSM Parameters")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}
	err = initTerraform(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing Terraform")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	//r := action.ActRunnerImpl{}
	//err = UpdateCoreConfig(appconf, store, nil, "")
	//if err != nil {
	//	log.WithError(err).Error("Problem updating ssm config")
	//}
	err = installGateway(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing HTTPD Gateway")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}
	
	err = installHealthStatusLambda(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing Health Status ")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	err = installBasicAPIGateway(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing API Gateway")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	err = store.AddToAudit(application.Bootstrap_Successful, "test")
	if err != nil {
		log.WithError(err).Error("Problem writing to auditlog")
	}

}

func provisionS3(appConfig *config.AppConfig) error {
	aws.CreateBucket(nil, appConfig)
	err := aws.CreateTable(appConfig, appConfig.InstallPrefix)
	if err != nil && !strings.Contains(err.Error(), "Table already exists") {
		log.WithError(err).Error("Error creating table")
		return err
	}

	return nil
}

func initTerraform(store database.Datastore, appconf *config.AppConfig) error {
	fs := afero.NewOsFs()
	err := writeInitTemplate(fs, appconf)
	if err != nil {
		return err
	}
	err = installUnityCloudEnv(store, appconf)
	if err != nil {
		return err
	}

	return nil

}

func writeInitTemplate(fs afero.Fs, appConfig *config.AppConfig) error {
	// Define the terraform configuration
	tfconfig := fmt.Sprintf(`terraform {
required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    dynamodb_table = "%s-terraform_state"
  }
}

provider "aws" {
  region = "us-west-2"
}`, appConfig.InstallPrefix)

	err := fs.MkdirAll(filepath.Join(appConfig.Workdir, "workspace"), 0755)
	if err != nil {
		log.WithError(err).Error("Couldn't create new working directory")
		return err
	}

	// Create a new file
	file, err := fs.Create(filepath.Join(appConfig.Workdir, "workspace", "provider.tf"))
	if err != nil {
		log.WithError(err).Error("Couldn't create new provider.tf file")
		return err
	}
	defer file.Close()

	// Write the configuration to the file
	_, err = file.WriteString(tfconfig)
	if err != nil {
		log.WithError(err).Error("Could not write provider string")
		return err
	}

	// Save changes to the file
	err = file.Sync()
	if err != nil {
		log.WithError(err).Error("Could not write provider.tf")
		return err
	}

	log.Println("File 'provider.tf' has been written")
	return nil
}
func storeDefaultSSMParameters(appConfig *config.AppConfig, store database.Datastore) error {

	err := store.StoreSSMParams(appConfig.DefaultSSMParameters, "bootstrap")
	if err != nil {
		log.WithError(err).Error("Problem storing parameters in database.")
		return err
	}
	return nil
}

func installGateway(store database.Datastore, appConfig *config.AppConfig) error {
	simplevars := make(map[string]string)
	simplevars["mgmt_dns"] = appConfig.ConsoleHost
	variables := marketplace.Install_Variables{Values: simplevars}
	applications := marketplace.Install_Applications{
		Name:        "unity-proxy",
		Version:     "0.14",
		Variables:   &variables,
		Displayname: fmt.Sprintf("%s-%s", appConfig.InstallPrefix, "unity-proxy"),
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Core Mgmt Gateway",
	}
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing Mgmt Gateway")
		return err
	}
	return nil
}

func installBasicAPIGateway(store database.Datastore, appConfig *config.AppConfig) error {
	applications := marketplace.Install_Applications{
		Name:        "unity-apigateway",
		Version:     "0.4",
		Variables:   nil,
		Displayname: fmt.Sprintf("%s-%s", appConfig.InstallPrefix, "unity-apigateway"),
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Core API Gateway",
	}
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing API Gateway")
		return err
	}
	return nil
}

func installUnityCloudEnv(store database.Datastore, appConfig *config.AppConfig) error {

	//venue, err := getSSMParameterValueFromDatabase("venue", store)
	//if err != nil {
	//	log.WithError(err).Error("Problem fetching venue")
	//	return err
	//}
	//log.Infof("Venue found: %s", venue)
	//project, err := getSSMParameterValueFromDatabase("project", store)
	//if err != nil {
	//	log.WithError(err).Error("Problem fetching project")
	//	return err
	//}
	//log.Infof("Project found: %s", project)
	project := appConfig.Project
	venue := appConfig.Venue

	if project == "" {
		log.Error("Config value Project not set")
	}
	if venue == "" {
		log.Error("Config value Venue not set")
	}

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
func installHealthStatusLambda(store database.Datastore, appConfig *config.AppConfig) error {
	applications := marketplace.Install_Applications{
		Name:        "unity-cs-monitoring-lambda",
		Version:     "0.1",
		Variables:   nil,
		Displayname: fmt.Sprintf("%s-%s", appConfig.InstallPrefix, "unity-cs-monitoring-lambda"),
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Unity Health Status Lambda",
	}
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing Unity Health Status Lambda")
		return err
	}
	return nil
}
