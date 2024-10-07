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

func InstallMarketplaceApplicationNew(appConfig *config.AppConfig, location string, installParams *types.ApplicationInstallParams, meta *marketplace.MarketplaceMetadata, db database.Datastore) (string, error) {
	if meta.Backend == "terraform" {
		app := models.Application{
			Name:        installParams.Name,
			Version:     installParams.Version,
			DisplayName: installParams.DeploymentName,
			PackageName: meta.Name,
			Source:      meta.Package,
			Status:      "STAGED",
		}
		deployment := models.Deployment{
			Name:         installParams.DeploymentName,
			Applications: []models.Application{app},
			Creator:      "admin",
			CreationDate: time.Time{},
		}

		deploymentID, err := db.StoreDeployment(deployment)

		if err != nil {
			db.UpdateApplicationStatus(deploymentID, installParams.Name, installParams.DeploymentName, "STAGINGFAILED")
			return "", err
		}

		go func() {
			log.Errorf("Application name is: %s", installParams.Name)
			err = terraform.AddApplicationToStackNew(appConfig, location, meta, installParams, db, deploymentID)
			executeNew(db, appConfig, meta, installParams, deploymentID)
		}()

		return fmt.Sprintf("%d", deploymentID), nil

	} else {
		return "", errors.New("backend not implemented")
	}
}

func InstallMarketplaceApplicationNewV2(appConfig *config.AppConfig, location string, installParams *types.ApplicationInstallParams, meta *marketplace.MarketplaceMetadata, db database.Datastore) error {
	if meta.Backend == "terraform" {
		application := models.InstalledMarketplaceApplication{
			Name:        installParams.Name,
			Version:     installParams.Version,
			DeploymentName: installParams.DeploymentName,
			PackageName: meta.Name,
			Source:      meta.Package,
			Status:      "STAGED",
		}

		db.StoreInstalledMarketplaceApplication(application)
		db.UpdateInstalledMarketplaceApplicationStatusByName(installParams.Name, installParams.DeploymentName, "INSTALLING")

		go func() {
			log.Errorf("Application name is: %s", installParams.Name)
			terraform.AddApplicationToStackNewV2(appConfig, location, meta, installParams, db)
			executeNewV2(db, appConfig, meta, installParams)
		}()

		return nil

	} else {
		return errors.New("backend not implemented")
	}
}

func InstallMarketplaceApplication(conn *websocket.WebSocketManager, userid string, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, location string, install *marketplace.Install, db database.Datastore) error {

	if meta.Backend == "terraform" {

		app := models.Application{
			Name:        install.Applications.Name,
			Version:     install.Applications.Version,
			DisplayName: install.Applications.Displayname,
			PackageName: meta.Name,
			Source:      meta.Package,
			Status:      "STAGED",
		}
		deployment := models.Deployment{
			Name:         install.DeploymentName,
			Applications: []models.Application{app},
			Creator:      "admin",
			CreationDate: time.Time{},
		}

		deploymentID, err := db.StoreDeployment(deployment)
		if err != nil {
			db.UpdateApplicationStatus(deploymentID, install.Applications.Name, install.Applications.Displayname, "STAGINGFAILED")

			return err
		}

		err = terraform.AddApplicationToStack(appConfig, location, meta, install, db, deploymentID)

		return execute(db, appConfig, meta, install, deploymentID, conn, userid)

	} else {
		return errors.New("backend not implemented")
	}
}

func execute(db database.Datastore, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, install *marketplace.Install, deploymentID uint, conn *websocket.WebSocketManager, userid string) error {
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
	err := runPreInstall(appConfig, meta, install)
	if err != nil {
		return err
	}
	db.UpdateApplicationStatus(deploymentID, install.Applications.Name, install.Applications.Displayname, "INSTALLING")
	fetchAllApplications(db)
	err = terraform.RunTerraform(appConfig, conn, userid, executor, "")
	if err != nil {
		db.UpdateApplicationStatus(deploymentID, install.Applications.Name, install.Applications.Displayname, "FAILED")
		fetchAllApplications(db)
		return err
	}
	db.UpdateApplicationStatus(deploymentID, install.Applications.Name, install.Applications.Displayname, "INSTALLED")
	fetchAllApplications(db)
	err = runPostInstall(appConfig, meta, install)

	if err != nil {
		db.UpdateApplicationStatus(deploymentID, install.Applications.Name, install.Applications.Displayname, "POSTINSTALL FAILED")
		fetchAllApplications(db)

		return err
	}
	db.UpdateApplicationStatus(deploymentID, install.Applications.Name, install.Applications.Displayname, "COMPLETE")
	fetchAllApplications(db)

	return nil
}

func executeNew(db database.Datastore, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, installParams *types.ApplicationInstallParams, deploymentID uint) error {
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
	err := runPreInstallNew(appConfig, meta, installParams)
	if err != nil {
		return err
	}

	db.UpdateApplicationStatus(deploymentID, installParams.Name, installParams.DeploymentName, "INSTALLING")

	fetchAllApplications(db)

	logfile := filepath.Join(logDir, fmt.Sprintf("%s_install_log", installParams.DeploymentName))
	err = terraform.RunTerraformLogOutToFile(appConfig, logfile, executor, "")

	if err != nil {
		db.UpdateApplicationStatus(deploymentID, installParams.Name, installParams.DeploymentName, "FAILED")
		fetchAllApplications(db)
		return err
	}
	db.UpdateApplicationStatus(deploymentID, installParams.Name, installParams.DeploymentName, "INSTALLED")
	fetchAllApplications(db)
	err = runPostInstallNew(appConfig, meta, installParams)

	if err != nil {
		db.UpdateApplicationStatus(deploymentID, installParams.Name, installParams.DeploymentName, "POSTINSTALL FAILED")
		fetchAllApplications(db)

		return err
	}
	db.UpdateApplicationStatus(deploymentID, installParams.Name, installParams.DeploymentName, "COMPLETE")
	fetchAllApplications(db)

	return nil
}

func executeNewV2(db database.Datastore, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, installParams *types.ApplicationInstallParams) error {
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
	err := runPreInstallNew(appConfig, meta, installParams)
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

func runPreInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, install *marketplace.Install) error {
	if meta.PreInstall != "" {
		// TODO UNPIN ME
		path := filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version, meta.WorkDirectory, meta.PreInstall)
		log.Infof("Pre install command path: %s", path)
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
			log.WithError(err).Error("Error running pre install script")
			return err
		}
	}
	return nil
}

func runPreInstallNew(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, installParams *types.ApplicationInstallParams) error {
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

func TriggerInstall(wsManager *websocket.WebSocketManager, userid string, store database.Datastore, received *marketplace.Install, conf *config.AppConfig) error {
	t := received.Applications

	meta, err := validateAndPrepareInstallation(t, conf)
	if err != nil {
		return err
	}

	location, err := FetchPackage(meta, conf)
	if err != nil {
		log.Error("Error fetching package:", err)
		return err
	}

	if err := InstallMarketplaceApplication(wsManager, userid, conf, meta, location, received, store); err != nil {
		log.Error("Error installing application:", err)
		return err
	}

	return nil
}

func TriggerUninstall(wsManager *websocket.WebSocketManager, userid string, store database.Datastore, received *marketplace.Uninstall, conf *config.AppConfig) error {
	if received.All == true {
		return UninstallAll(conf, wsManager, userid, received)
	} else {
		return UninstallApplication(received.Application, received.DeploymentName, received.DisplayName, conf, store, wsManager, userid)
	}
}
