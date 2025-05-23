package processes

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	// "github.com/unity-sds/unity-cs-manager/marketplace"

	"github.com/unity-sds/unity-management-console/backend/internal/application"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/types"
)

func BootstrapEnv(appconf *config.AppConfig) error {
	// Print out everything in appConfig
	log.Infof("AppConfig contents:")
	log.Infof("MarketplaceBaseUrl: %s", appconf.MarketplaceBaseUrl)
	log.Infof("MarketplaceOwner: %s", appconf.MarketplaceOwner)
	log.Infof("MarketplaceRepo: %s", appconf.MarketplaceRepo)
	log.Infof("AWSRegion: %s", appconf.AWSRegion)
	log.Infof("BucketName: %s", appconf.BucketName)
	log.Infof("Workdir: %s", appconf.Workdir)
	log.Infof("BasePath: %s", appconf.BasePath)
	log.Infof("ConsoleHost: %s", appconf.ConsoleHost)
	log.Infof("InstallPrefix: %s", appconf.InstallPrefix)
	log.Infof("Project: %s", appconf.Project)
	log.Infof("Venue: %s", appconf.Venue)
	log.Infof("MarketplaceItems:")
	for _, item := range appconf.MarketplaceItems {
		log.Infof("  - Name: %s, Version: %s", item.Name, item.Version)
	}

	log.Infof("Creating Local Database")
	store, err := database.NewGormDatastore()
	if err != nil {
		log.WithError(err).Error("Problem creating database")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return err
	}

	log.Infof("Provisioning S3 Bucket")
	err = provisionS3(appconf)
	if err != nil {
		log.WithError(err).Error("Error provisioning S3 bucket")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return err
	}

	log.Infof("Setting Up Default SSM Parameters")
	err = storeDefaultSSMParameters(appconf, store)
	if err != nil {
		log.WithError(err).Error("Error setting SSM Parameters")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return err
	}

	log.Infof("Setting Up Terraform")
	err = initTerraform(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing Terraform")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return err
	}

	//r := action.ActRunnerImpl{}
	//err = UpdateCoreConfig(appconf, store, nil, "")
	//if err != nil {
	//	log.WithError(err).Error("Problem updating ssm config")
	//}

	log.Infof("Setting Up HTTPD Gateway and API Gateway from Marketplace in parallel")
	err = installGatewayAndApiGateway(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing Gateways")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return err
	}

	log.Infof("Setting Up Health Status Lambda and Unity UI from Marketplace in parallel")
	err = installHealthStatusLambdaAndUnityUi(store, appconf)
	if err != nil {
		log.WithError(err).Error("Error installing Health Status Lambda and Unity UI")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return err
	}

	err = store.AddToAudit(application.Bootstrap_Successful, "test")
	if err != nil {
		log.WithError(err).Error("Problem writing to auditlog")
		return err
	}

	log.Infof("Setting Up Core Config SSM Params")
	err = UpdateCoreConfig(appconf, store, nil, "")
	if err != nil {
		log.WithError(err).Error("Problem updating ssm config")
	}
	store.AddToAudit(application.Config_Updated, "test")
	
	log.Infof("Bootstrap Process Completed Succesfully")
	return nil
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

	// err = installUnityCloudEnv(store, appconf)
	// if err != nil {
	// 	return err
	// }

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
    use_lockfile = true 
  }
}

provider "aws" {
  region = "us-west-2"
}
`)

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

func installGatewayAndApiGateway(store database.Datastore, appConfig *config.AppConfig) error {
	// Find the marketplace item for unity-proxy
	var proxyName, proxyVersion string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-proxy" {
			proxyName = item.Name
			proxyVersion = item.Version
			break
		}
	}
	if proxyName == "" || proxyVersion == "" {
		log.Error("unity-proxy not found in MarketplaceItems")
		return fmt.Errorf("unity-proxy not found in MarketplaceItems")
	}
	
	// Find the marketplace item for unity-apigateway
	var apiName, apiVersion string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-apigateway" {
			apiName = item.Name
			apiVersion = item.Version
			break
		}
	}
	if apiName == "" || apiVersion == "" {
		log.Error("unity-apigateway not found in MarketplaceItems")
		return fmt.Errorf("unity-apigateway not found in MarketplaceItems")
	}
	
	// Set up variables for unity-proxy
	proxyVars := make(map[string]string)
	proxyVars["mgmt_dns"] = appConfig.ConsoleHost
	
	// Create installation parameters for both applications
	params := []*types.ApplicationInstallParams{
		{
			Name:           proxyName,
			Version:        proxyVersion,
			Variables:      proxyVars,
			DisplayName:    "Core Mgmt Gateway",
			DeploymentName: fmt.Sprintf("default-%s", proxyName),
		},
		{
			Name:           apiName,
			Version:        apiVersion,
			Variables:      nil,
			DisplayName:    "Core API Gateway",
			DeploymentName: fmt.Sprintf("default-%s", apiName),
		},
	}
	
	// Install both applications in a single batch operation
	return BatchTriggerInstall(store, params, appConfig)
}

func installHealthStatusLambdaAndUnityUi(store database.Datastore, appConfig *config.AppConfig) error {
	// Find the marketplace item for health status lambda
	var lambdaName, lambdaVersion string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-cs-monitoring-lambda" {
			lambdaName = item.Name
			lambdaVersion = item.Version
			break
		}
	}
	if lambdaName == "" || lambdaVersion == "" {
		log.Error("unity-cs-monitoring-lambda not found in MarketplaceItems")
		return fmt.Errorf("unity-cs-monitoring-lambda not found in MarketplaceItems")
	}
	
	// Find the marketplace item for unity-portal
	var uiName, uiVersion string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-portal" {
			uiName = item.Name
			uiVersion = item.Version
			break
		}
	}
	if uiName == "" || uiVersion == "" {
		log.Error("unity-portal not found in MarketplaceItems")
		return fmt.Errorf("unity-portal not found in MarketplaceItems")
	}
	
	// Create installation parameters for both applications
	params := []*types.ApplicationInstallParams{
		{
			Name:           lambdaName,
			Version:        lambdaVersion,
			Variables:      nil,
			DisplayName:    "Unity Health Status Lambda",
			DeploymentName: fmt.Sprintf("default-%s", lambdaName),
		},
		{
			Name:           uiName,
			Version:        uiVersion,
			Variables:      nil,
			DisplayName:    "Unity Navbar UI",
			DeploymentName: fmt.Sprintf("default-%s", uiName),
		},
	}
	
	// Install both applications in a single batch operation
	return BatchTriggerInstall(store, params, appConfig)
}

func installUnityCloudEnv(store database.Datastore, appConfig *config.AppConfig) error {
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

	// Find the marketplace item for unity-cloud-env
	var name, version string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-cloud-env" {
			name = item.Name
			version = item.Version
			break
		}
	}

	// Print the name and version
	log.Infof("Found marketplace item - Name: %s, Version: %s", name, version)

	// If the item wasn't found, log an error and return
	if name == "" || version == "" {
		log.Error("unity-cloud-env not found in MarketplaceItems")
		return fmt.Errorf("unity-cloud-env not found in MarketplaceItems")
	}

	varmap := make(map[string]string)
	varmap["venue"] = venue
	varmap["project"] = project
	varmap["publicsubnets"] = publicsubnets
	varmap["privatesubnets"] = privatesubnets
	// vars := marketplace.Install_Variables{
	// 	Values:         varmap,
	// 	AdvancedValues: nil,
	// }
	// applications := marketplace.Install_Applications{
	// 	Name:        name,
	// 	Version:     version,
	// 	Variables:   &vars,
	// 	Displayname: name,
	// }
	// install := marketplace.Install{
	// 	Applications:   &applications,
	// 	DeploymentName: "Unity Cloud Environment",
	// }

	installParams := types.ApplicationInstallParams{
		Name:           name,
		Version:        version,
		Variables:      varmap,
		DisplayName:    "Unity Cloud Environment",
		DeploymentName: fmt.Sprintf("default-%s", name),
	}

	err = TriggerInstall(store, &installParams, appConfig, true)
	if err != nil {
		log.WithError(err).Error("Issue installing Unity Cloud Env")
		return err
	}
	return nil
}
