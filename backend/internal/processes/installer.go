package processes

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"github.com/unity-sds/unity-management-console/backend/types"
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	// "time"
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



func startApplicationInstallTerraform(appConfig *config.AppConfig, location string, installParams *types.ApplicationInstallParams, meta *marketplace.MarketplaceMetadata, db database.Datastore) {
	log.Errorf("Application name is: %s", installParams.Name)
	terraform.AddApplicationToStack(appConfig, location, meta, installParams, db)
	execute(db, appConfig, meta, installParams)
}

func InstallMarketplaceApplication(appConfig *config.AppConfig, location string, installParams *types.ApplicationInstallParams, meta *marketplace.MarketplaceMetadata, db database.Datastore, sync bool) error {
	if meta.Backend == "terraform" {
		application := models.InstalledMarketplaceApplication{
			Name:           installParams.Name,
			Version:        installParams.Version,
			DeploymentName: installParams.DeploymentName,
			PackageName:    meta.Name,
			Source:         meta.Package,
			Status:         "STAGED",
		}

		db.StoreInstalledMarketplaceApplication(application)
		db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "INSTALLING")

		if sync {
			startApplicationInstallTerraform(appConfig, location, installParams, meta, db)
		} else {
			go startApplicationInstallTerraform(appConfig, location, installParams, meta, db)

		}

		return nil

	} else {
		return errors.New("backend not implemented")
	}
}





func execute(db database.Datastore, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, installParams *types.ApplicationInstallParams) error {
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
	err := runPreInstall(appConfig, meta, installParams)
	if err != nil {
		return err
	}

	db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "INSTALLING")

	fetchAllApplications(db)

	logfile := filepath.Join(logDir, fmt.Sprintf("%s_%s_install_log", installParams.Name, installParams.DeploymentName))
	err = terraform.RunTerraformLogOutToFile(appConfig, logfile, executor, "")

	if err != nil {
		db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "FAILED")
		fetchAllApplications(db)
		return err
	}
	db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "INSTALLED")
	fetchAllApplications(db)
	err = runPostInstallNew(appConfig, meta, installParams)

	if err != nil {
		db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "POSTINSTALL FAILED")
		fetchAllApplications(db)

		return err
	}
	db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "COMPLETE")
	fetchAllApplications(db)

	return nil
}
func runPostInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, install *marketplace.Install) error {

	if meta.PostInstall != "" {
		//TODO UNPIN ME
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.WorkDirectory, meta.PostInstall)
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

func runPostInstallNew(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, installParams *types.ApplicationInstallParams) error {

	if meta.PostInstall != "" {
		//TODO UNPIN ME
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.WorkDirectory, meta.PostInstall)
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
		cmd.Env = append(cmd.Env, fmt.Sprintf("DEPLOYMENTNAME=%s", installParams.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", appConfig.Workdir))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DISPLAYNAME=%s", installParams.DisplayName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("APPNAME=%s", installParams.Name))

		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running post install script")
			return err
		}
	}
	return nil
}


func runPreInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, installParams *types.ApplicationInstallParams) error {
	if meta.PreInstall != "" {
		// TODO UNPIN ME
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.WorkDirectory, meta.PreInstall)
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
		cmd.Env = append(cmd.Env, fmt.Sprintf("DEPLOYMENTNAME=%s", installParams.DeploymentName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", appConfig.Workdir))
		cmd.Env = append(cmd.Env, fmt.Sprintf("DISPLAYNAME=%s", installParams.DisplayName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("APPNAME=%s", installParams.Name))
		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running pre install script")
			return err
		}
	}
	return nil
}


func TriggerInstall(store database.Datastore, applicationInstallParams *types.ApplicationInstallParams, conf *config.AppConfig, sync bool) error {
	// First check if this application is already installed.
	existingApplication, err := store.GetInstalledMarketplaceApplicationStatusByName(applicationInstallParams.Name, applicationInstallParams.DeploymentName)
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

func TriggerUninstall(wsManager *websocket.WebSocketManager, userid string, store database.Datastore, received *marketplace.Uninstall, conf *config.AppConfig) error {
	if received.All == true {
		return UninstallAll(conf, wsManager, userid, received)
	} else {
		return UninstallApplication(received.Application, received.DeploymentName, received.DisplayName, conf, store, wsManager, userid)
	}
}
