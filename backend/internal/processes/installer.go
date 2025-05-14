package processes

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"github.com/unity-sds/unity-management-console/backend/types"
	"github.com/zclconf/go-cty/cty"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func fetchMandatoryVars() ([]terraform.Varstruct, error) {
	pub, priv, err := aws.FetchSubnets()
	if err != nil {
		return nil, err
	}
	varType := cty.Map(cty.List(cty.String))
	vars := []terraform.Varstruct{
		{"vpc_id", "vpc-03dfddb7c55b63591", cty.String},
		{"subnets", map[string][]string{"public": pub, "private": priv}, varType},
	}

	return vars, nil

}

func startApplicationInstallTerraform(appConfig *config.AppConfig, location string, application *types.InstalledMarketplaceApplication, meta *types.MarketplaceMetadata, db database.Datastore) {
	log.Errorf("Application name is: %s", application.Name)
	terraform.AddApplicationToStack(appConfig, location, meta, application, db)
	executeTerraformInstall(db, appConfig, meta, application)
}

func InstallMarketplaceApplication(appConfig *config.AppConfig, location string, installParams *types.ApplicationInstallParams, meta *types.MarketplaceMetadata, db database.Datastore, sync bool) error {
	if meta.TerraformModuleName != "" {

		rand.Seed(time.Now().UnixNano())
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		randomChars := make([]byte, 5)
		for i, v := range rand.Perm(52)[:5] {
			randomChars[i] = chars[v]
		}
		log.Info("Appending block to body")
		terraformModuleName := fmt.Sprintf("%s-%s", installParams.DeploymentName, string(randomChars))

		application := &types.InstalledMarketplaceApplication{
			Name:                installParams.Name,
			Version:             installParams.Version,
			DeploymentName:      installParams.DeploymentName,
			PackageName:         meta.Name,
			Source:              meta.Package,
			Status:              "STAGED",
			TerraformModuleName: terraformModuleName,
			Variables:           installParams.Variables,
		}

		db.StoreInstalledMarketplaceApplication(application)

		application.Status = "INSTALLING"
		db.UpdateInstalledMarketplaceApplication(application)

		if sync {
			startApplicationInstallTerraform(appConfig, location, application, meta, db)
		} else {
			go startApplicationInstallTerraform(appConfig, location, application, meta, db)

		}

		return nil

	} else {
		return errors.New("backend not implemented")
	}
}

func executeTerraformInstall(db database.Datastore, appConfig *config.AppConfig, meta *types.MarketplaceMetadata, application *types.InstalledMarketplaceApplication) error {
	// Create install_logs directory if it doesn't exist
	logDir := filepath.Join(appConfig.Workdir, "install_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create install_logs directory: %w", err)
	}

	executor := &terraform.RealTerraformExecutor{}

	//m, err := fetchMandatoryVars()
	//if err != nil {
	//	return err
	//}
	//terraform.WriteTFVars(m, appConfig)
	err := runPreInstall(appConfig, meta, application)
	if err != nil {
		return err
	}

	application.Status = "INSTALLING"
	db.UpdateInstalledMarketplaceApplication(application)

	fetchAllApplications(db)

	logfile := filepath.Join(logDir, fmt.Sprintf("%s_%s_install_log", application.Name, application.DeploymentName))
	err = terraform.RunTerraformLogOutToFile(appConfig, logfile, executor, "")

	if err != nil {
		application.Status = "FAILED"
		db.UpdateInstalledMarketplaceApplication(application)
		return err
	}

	application.Status = "INSTALLED"
	db.UpdateInstalledMarketplaceApplication(application)
	fetchAllApplications(db)
	err = runPostInstallNew(appConfig, meta, application)

	if err != nil {
		application.Status = "POSTINSTALL FAILED"
		db.UpdateInstalledMarketplaceApplication(application)
		fetchAllApplications(db)
		return err
	}

	// Check to make sure the required output SSM params have been populated
	results, err := checkInstalledDependencies(*meta, appConfig)
	if err != nil {
		application.Status = "OUTPUT SSM PARAM CHECK FAILED"
		db.UpdateInstalledMarketplaceApplication(application)
		return fmt.Errorf("Error getting required SSM params: %v", results)
	}

	application.Status = "COMPLETE"
	db.UpdateInstalledMarketplaceApplication(application)
	fetchAllApplications(db)

	return nil
}
func runPostInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, install *marketplace.Install) error {

	if meta.PostInstall != "" {
		//TODO UNPIN ME
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.PostInstall)
		log.Infof("Post install command path: %s", path)
		cmd := exec.Command(path)
		cmd.Env = os.Environ()
		for k, v := range install.Applications.Dependencies {
			// Replace hyphens with underscores
			formattedKey := strings.ReplaceAll(k, "-", "_")

			// Convert to upper case
			upperKey := strings.ToUpper(formattedKey)

			// Use a regex to keep only alphanumeric characters and underscores
			re := regexp.MustCompile("[^A-Z0-9_]+")
			cleanKey := re.ReplaceAllString(upperKey, "")

			log.Infof("Adding environment variable: %s = %s", cleanKey, v)
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", cleanKey, v))
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("DEPLOYMENTNAME=%s", install.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", appConfig.Workdir))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DISPLAYNAME=%s", install.Applications.Displayname))
		cmd.Env = append(cmd.Env, fmt.Sprintf("APPNAME=%s", install.Applications.Name))

		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running post install script")
			return err
		}
	}
	return nil
}

func runPostInstallNew(appConfig *config.AppConfig, meta *types.MarketplaceMetadata, application *types.InstalledMarketplaceApplication) error {

	if meta.PostInstall != "" {
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.PostInstall)
		log.Infof("Post install command path: %s", path)
		cmd := exec.Command(path)
		cmd.Env = os.Environ()
		// for k, v := range install.Applications.Dependencies {
		// 	// Replace hyphens with underscores
		// 	formattedKey := strings.ReplaceAll(k, "-", "_")

		// 	// Convert to upper case
		// 	upperKey := strings.ToUpper(formattedKey)

		// 	// Use a regex to keep only alphanumeric characters and underscores
		// 	re := regexp.MustCompile("[^A-Z0-9_]+")
		// 	cleanKey := re.ReplaceAllString(upperKey, "")

		// 	log.Infof("Adding environment variable: %s = %s", cleanKey, v)
		// 	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", cleanKey, v))
		// }
		cmd.Env = append(cmd.Env, fmt.Sprintf("DEPLOYMENTNAME=%s", application.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", appConfig.Workdir))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DISPLAYNAME=%s", application.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("APPNAME=%s", application.Name))

		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running post install script")
			return err
		}
	}
	return nil
}

func runPreInstall(appConfig *config.AppConfig, meta *types.MarketplaceMetadata, application *types.InstalledMarketplaceApplication) error {
	if meta.PreInstall != "" {
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.PreInstall)
		log.Infof("Pre install command path: %s", path)
		cmd := exec.Command(path)
		cmd.Env = os.Environ()
		// for k, v := range install.Applications.Dependencies {
		// 	// Replace hyphens with underscores
		// 	formattedKey := strings.ReplaceAll(k, "-", "_")

		// 	// Convert to upper case
		// 	upperKey := strings.ToUpper(formattedKey)

		// 	// Use a regex to keep only alphanumeric characters and underscores
		// 	re := regexp.MustCompile("[^A-Z0-9_]+")
		// 	cleanKey := re.ReplaceAllString(upperKey, "")

		// 	log.Infof("Adding environment variable: %s = %s", cleanKey, v)
		// 	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", cleanKey, v))
		// }
		cmd.Env = append(cmd.Env, fmt.Sprintf("DEPLOYMENTNAME=%s", application.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", appConfig.Workdir))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DISPLAYNAME=%s", application.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("APPNAME=%s", application.Name))
		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running pre install script")
			return err
		}
	}
	return nil
}

func TriggerInstall(store database.Datastore, applicationInstallParams *types.ApplicationInstallParams, conf *config.AppConfig, sync bool) error {
	// First check if this application is already installed.
	existingApplication, err := store.GetInstalledMarketplaceApplication(applicationInstallParams.Name, applicationInstallParams.DeploymentName)
	if err != nil {
		log.WithError(err).Error("Error finding applications")
		return errors.New("Unable to search applcation list")
	}

	if existingApplication != nil && existingApplication.Status != "UNINSTALLED" {
		errMsg := fmt.Sprintf("Application with name %s already exists. Status: %s", applicationInstallParams.Name, existingApplication.Status)
		return errors.New(errMsg)
	}

	// OK to start installing, get the metadata for this application
	metadata, err := FetchMarketplaceMetadata(applicationInstallParams.Name, applicationInstallParams.Version, conf)
	if err != nil {
		log.Errorf("Unable to fetch metadata for application: %s, %v", applicationInstallParams.Name, err)
		return errors.New("Unable to fetch package")
	}

	// Install the application package files
	location, err := FetchPackage(&metadata, conf)
	if err != nil {
		log.Errorf("Unable to fetch package for application: %s, %v", applicationInstallParams.Name, err)
		return errors.New("Unable to fetch package")
	}

	return InstallMarketplaceApplication(conf, location, applicationInstallParams, &metadata, store, sync)
}

// BatchTriggerInstall handles the installation of multiple applications in a single Terraform operation
// to avoid state locking issues when installing multiple modules in parallel
func BatchTriggerInstall(store database.Datastore, params []*types.ApplicationInstallParams, conf *config.AppConfig) error {
	if len(params) == 0 {
		return nil
	}

	applications := make([]*types.InstalledMarketplaceApplication, 0, len(params))
	metadatas := make([]*types.MarketplaceMetadata, 0, len(params))
	locations := make([]string, 0, len(params))

	// First pass: validate and prepare all applications
	for _, param := range params {
		// Check if application is already installed
		existingApplication, err := store.GetInstalledMarketplaceApplication(param.Name, param.DeploymentName)
		if err != nil {
			log.WithError(err).Error("Error finding applications")
			return errors.New("Unable to search application list")
		}

		if existingApplication != nil && existingApplication.Status != "UNINSTALLED" {
			errMsg := fmt.Sprintf("Application with name %s already exists. Status: %s", param.Name, existingApplication.Status)
			return errors.New(errMsg)
		}

		// Fetch metadata and package
		metadata, err := FetchMarketplaceMetadata(param.Name, param.Version, conf)
		if err != nil {
			log.Errorf("Unable to fetch metadata for application: %s, %v", param.Name, err)
			return errors.New("Unable to fetch package")
		}

		location, err := FetchPackage(&metadata, conf)
		if err != nil {
			log.Errorf("Unable to fetch package for application: %s, %v", param.Name, err)
			return errors.New("Unable to fetch package")
		}

		// Generate unique module name
		rand.Seed(time.Now().UnixNano())
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		randomChars := make([]byte, 5)
		for i, v := range rand.Perm(52)[:5] {
			randomChars[i] = chars[v]
		}
		terraformModuleName := fmt.Sprintf("%s-%s", param.DeploymentName, string(randomChars))

		// Create application record
		application := &types.InstalledMarketplaceApplication{
			Name:                param.Name,
			Version:             param.Version,
			DeploymentName:      param.DeploymentName,
			PackageName:         metadata.Name,
			Source:              metadata.Package,
			Status:              "STAGED",
			TerraformModuleName: terraformModuleName,
			Variables:           param.Variables,
		}

		// Store application record
		store.StoreInstalledMarketplaceApplication(application)

		// Update status to INSTALLING
		application.Status = "INSTALLING"
		store.UpdateInstalledMarketplaceApplication(application)

		// Add to our tracking lists
		applications = append(applications, application)
		metadatas = append(metadatas, &metadata)
		locations = append(locations, location)
	}

	// Second pass: prepare all Terraform configurations
	for i, application := range applications {
		// Run pre-install script if needed
		err := runPreInstall(conf, metadatas[i], application)
		if err != nil {
			log.WithError(err).Errorf("Error running pre-install for %s", application.Name)
			application.Status = "FAILED"
			store.UpdateInstalledMarketplaceApplication(application)
			return err
		}

		// Add application to Terraform stack
		err = terraform.AddApplicationToStack(conf, locations[i], metadatas[i], application, store)
		if err != nil {
			log.WithError(err).Errorf("Error adding application to stack: %s", application.Name)
			application.Status = "FAILED"
			store.UpdateInstalledMarketplaceApplication(application)
			return err
		}
	}

	// Third pass: run a single Terraform apply for all modules
	executor := &terraform.RealTerraformExecutor{}
	logDir := filepath.Join(conf.Workdir, "install_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create install_logs directory: %w", err)
	}

	logfile := filepath.Join(logDir, "batch_install_log")
	err := terraform.RunTerraformLogOutToFile(conf, logfile, executor, "")
	if err != nil {
		// Mark all applications as failed
		for _, application := range applications {
			application.Status = "FAILED"
			store.UpdateInstalledMarketplaceApplication(application)
		}
		return err
	}

	// Fourth pass: run post-install scripts and update status
	for i, application := range applications {
		application.Status = "INSTALLED"
		store.UpdateInstalledMarketplaceApplication(application)

		err := runPostInstallNew(conf, metadatas[i], application)
		if err != nil {
			log.WithError(err).Errorf("Error running post-install for %s", application.Name)
			application.Status = "POSTINSTALL FAILED"
			store.UpdateInstalledMarketplaceApplication(application)
			return err
		}

		application.Status = "COMPLETE"
		store.UpdateInstalledMarketplaceApplication(application)
	}

	return nil
}

func TriggerUninstall(wsManager *websocket.WebSocketManager, userid string, store database.Datastore, received *marketplace.Uninstall, conf *config.AppConfig) error {
	if received.All == true {
		return UninstallAll(conf, wsManager, userid, received)
	}
	// else {
	// 	return UninstallApplication(received.Application, received.DeploymentName, received.DisplayName, conf, store, wsManager, userid)
	// }
	return nil
}

func checkInstalledDependencies(metadata types.MarketplaceMetadata, conf *config.AppConfig) (map[string]string, error) {
	errorsEncountered := false
	results := make(map[string]string)
	for _, ssmParamName := range metadata.OutputSsmParameters {
		formattedParam := strings.Replace(ssmParamName, "${PROJ}", conf.Project, -1)
		formattedParam = strings.Replace(formattedParam, "${VENUE}", conf.Venue, -1)

		log.Infof("Looking up key %s", formattedParam)
		param, err := aws.ReadSSMParameter(formattedParam)

		if err != nil {
			log.WithError(err).Error("Unable to get SSM param.")
			results[ssmParamName] = ""
			errorsEncountered = true
			continue
		}
		log.Infof("Got param %s", *param.Parameter.Value)
		results[ssmParamName] = *param.Parameter.Value
	}

	if errorsEncountered {
		return results, fmt.Errorf("Output SSM Params Not found")
	}

	return results, nil
}
