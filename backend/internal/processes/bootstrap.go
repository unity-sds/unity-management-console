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

// BootstrapEnv initializes the environment for the Unity Management Console.
// It sets up necessary AWS resources, installs core components, and performs initial configurations.
func BootstrapEnv(appconf *config.AppConfig) {
	// Log all AppConfig contents for debugging purposes
	// This helps in verifying that all configuration parameters are correctly loaded
	log.Infof("AppConfig contents:")
	// Print out everything in appConfig
    log.Infof("GithubToken: %s", appconf.GithubToken)
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

	// Initialize the database connection using GORM
	// This creates a new datastore instance for interacting with the database
	store, err := database.NewGormDatastore()
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt in the audit log
		log.WithError(err).Error("Problem creating database")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Set up S3 bucket and DynamoDB table for Terraform state management
	// This ensures that Terraform has a place to store its state files
	err = provisionS3(appconf)
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt
		log.WithError(err).Error("Error provisioning S3 bucket")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Store default SSM (Systems Manager) parameters in the database
	// These parameters are used for configuration management across the system
	err = storeDefaultSSMParameters(appconf, store)
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt
		log.WithError(err).Error("Error setting SSM Parameters")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Initialize Terraform configuration
	// This sets up the basic Terraform files needed for managing infrastructure
	err = initTerraform(store, appconf)
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt
		log.WithError(err).Error("Error installing Terraform")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Install the HTTP Gateway component
	// This sets up the main entry point for HTTP traffic to the Unity system
	err = installGateway(store, appconf)
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt
		log.WithError(err).Error("Error installing HTTPD Gateway")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Install the Health Status Lambda function
	// This function is used to monitor and report on the health of the Unity system
	err = installHealthStatusLambda(store, appconf)
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt
		log.WithError(err).Error("Error installing Health Status")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Install the basic API Gateway
	// This sets up the API Gateway for routing requests to various Unity services
	err = installBasicAPIGateway(store, appconf)
	if err != nil {
		// Log error and record unsuccessful bootstrap attempt
		log.WithError(err).Error("Error installing API Gateway")
		err = store.AddToAudit(application.Bootstrap_Unsuccessful, "test")
		if err != nil {
			log.WithError(err).Error("Problem writing to auditlog")
		}
		return
	}

	// Record successful bootstrap completion in the audit log
	// This helps track the successful setup of the Unity environment
	err = store.AddToAudit(application.Bootstrap_Successful, "test")
	if err != nil {
		log.WithError(err).Error("Problem writing to auditlog")
	}

}

// provisionS3 creates an S3 bucket and DynamoDB table for Terraform state management
func provisionS3(appConfig *config.AppConfig) error {
	// Create S3 bucket for storing Terraform state files
	// This bucket will be used to store the state of all Terraform-managed resources
	aws.CreateBucket(nil, appConfig)

	// Create DynamoDB table for Terraform state locking
	// This table prevents concurrent modifications to the Terraform state
	err := aws.CreateTable(appConfig, appConfig.InstallPrefix)
	if err != nil && !strings.Contains(err.Error(), "Table already exists") {
		log.WithError(err).Error("Error creating table")
		return err
	}

	return nil
}

// initTerraform sets up the initial Terraform configuration
func initTerraform(store database.Datastore, appconf *config.AppConfig) error {
	// Create a new filesystem interface for file operations
	fs := afero.NewOsFs()

	// Write the initial Terraform configuration template
	// This creates the basic provider.tf file with AWS provider settings
	err := writeInitTemplate(fs, appconf)
	if err != nil {
		return err
	}

	// Note: The following code is commented out, but it's worth mentioning
	// that it would install the Unity Cloud Environment if uncommented
	// err = installUnityCloudEnv(store, appconf)
	// if err != nil {
	// 	return err
	// }

	return nil

}

// writeInitTemplate creates the initial Terraform configuration file
func writeInitTemplate(fs afero.Fs, appConfig *config.AppConfig) error {
	// Define the Terraform configuration content
	// This includes AWS provider setup and backend configuration for state management
	tfconfig := fmt.Sprintf(`terraform {
required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    dynamodb_table = "%s-%s-terraform-state"
  }
}

provider "aws" {
  region = "us-west-2"
}`, appConfig.Project, appConfig.Venue)

	err := fs.MkdirAll(filepath.Join(appConfig.Workdir, "workspace"), 0755)
	if err != nil {
		log.WithError(err).Error("Couldn't create new working directory")
		return err
	}

	// Create and write to the provider.tf file
	// This file contains the main Terraform configuration for the AWS provider
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

// storeDefaultSSMParameters stores the default SSM parameters in the database
func storeDefaultSSMParameters(appConfig *config.AppConfig, store database.Datastore) error {
	// Store the default SSM parameters in the database
	// These parameters are used for various configuration settings across the system
	err := store.StoreSSMParams(appConfig.DefaultSSMParameters, "bootstrap")
	if err != nil {
		log.WithError(err).Error("Problem storing parameters in database.")
		return err
	}
	return nil
}

// installGateway installs the Unity Proxy (HTTP Gateway) component
func installGateway(store database.Datastore, appConfig *config.AppConfig) error {
	// Find the marketplace item for unity-proxy
	// This component acts as the main HTTP gateway for the Unity system
	var name, version string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-proxy" {
			name = item.Name
			version = item.Version
			break
		}
	}

	// Print the name and version for logging purposes
	log.Infof("Found marketplace item - Name: %s, Version: %s", name, version)

	// If the item wasn't found, log an error and return
	if name == "" || version == "" {
		log.Error("unity-proxy not found in MarketplaceItems")
		return fmt.Errorf("unity-proxy not found in MarketplaceItems")
	}

	// Prepare the installation parameters
	simplevars := make(map[string]string)
	simplevars["mgmt_dns"] = appConfig.ConsoleHost
	variables := marketplace.Install_Variables{Values: simplevars}
	applications := marketplace.Install_Applications{
		Name:        name,
		Version:     version,
		Variables:   &variables,
		Displayname: fmt.Sprintf("%s-%s", appConfig.InstallPrefix, name),
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Core Mgmt Gateway",
	}

	// Trigger the installation process
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing Mgmt Gateway")
		return err
	}
	return nil
}

// installBasicAPIGateway installs the basic API Gateway component
func installBasicAPIGateway(store database.Datastore, appConfig *config.AppConfig) error {
	// Find the marketplace item for unity-apigateway
	// This component sets up the API Gateway for routing requests to Unity services
	var name, version string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-apigateway" {
			name = item.Name
			version = item.Version
			break
		}
	}

	// Print the name and version for logging purposes
	log.Infof("Found marketplace item - Name: %s, Version: %s", name, version)

	// If the item wasn't found, log an error and return
	if name == "" || version == "" {
		log.Error("unity-apigateway not found in MarketplaceItems")
		return fmt.Errorf("unity-apigateway not found in MarketplaceItems")
	}

	// Prepare the installation parameters
	applications := marketplace.Install_Applications{
		Name:        name,
		Version:     version,
		Variables:   nil,
		Displayname: fmt.Sprintf("%s-%s", appConfig.InstallPrefix, name),
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Core API Gateway",
	}

	// Trigger the installation process
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing API Gateway")
		return err
	}
	return nil
}

// installUnityCloudEnv installs the Unity Cloud Environment component
func installUnityCloudEnv(store database.Datastore, appConfig *config.AppConfig) error {
	// Retrieve project and venue information from the configuration
	project := appConfig.Project
	venue := appConfig.Venue

	// Validate that project and venue are set
	if project == "" {
		log.Error("Config value Project not set")
	}
	if venue == "" {
		log.Error("Config value Venue not set")
	}

	// Fetch public and private subnet information from SSM parameters
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

	// Print the name and version for logging purposes
	log.Infof("Found marketplace item - Name: %s, Version: %s", name, version)

	// If the item wasn't found, log an error and return
	if name == "" || version == "" {
		log.Error("unity-cloud-env not found in MarketplaceItems")
		return fmt.Errorf("unity-cloud-env not found in MarketplaceItems")
	}

	// Prepare the installation parameters
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
		Name:        name,
		Version:     version,
		Variables:   &vars,
		Displayname: name,
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Unity Cloud Environment",
	}

	// Trigger the installation process
	err = TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing Unity Cloud Env")
		return err
	}
	return nil
}

// installHealthStatusLambda installs the Health Status Lambda function
func installHealthStatusLambda(store database.Datastore, appConfig *config.AppConfig) error {
	// Find the marketplace item for the health status lambda
	// This Lambda function is responsible for monitoring and reporting system health
	var name, version string
	for _, item := range appConfig.MarketplaceItems {
		if item.Name == "unity-cs-monitoring-lambda" {
			name = item.Name
			version = item.Version
			break
		}
	}

	// Print the name and version for logging purposes
	log.Infof("Found marketplace item - Name: %s, Version: %s", name, version)

	// If the item wasn't found, log an error and return
	if name == "" || version == "" {
		log.Error("unity-cs-monitoring-lambda not found in MarketplaceItems")
		return fmt.Errorf("unity-cs-monitoring-lambda not found in MarketplaceItems")
	}

	// Prepare the installation parameters
	applications := marketplace.Install_Applications{
		Name:        name,
		Version:     version,
		Variables:   nil,
		Displayname: fmt.Sprintf("%s-%s", appConfig.InstallPrefix, name),
	}
	install := marketplace.Install{
		Applications:   &applications,
		DeploymentName: "Unity Health Status Lambda",
	}

	// Trigger the installation process
	err := TriggerInstall(nil, "", store, &install, appConfig)
	if err != nil {
		log.WithError(err).Error("Issue installing Unity Health Status Lambda")
		return err
	}
	return nil
}