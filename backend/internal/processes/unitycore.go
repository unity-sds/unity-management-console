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
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"path/filepath"
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
func InstallMarketplaceApplication(conn *websocket.WebSocketManager, userid string, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, location string, install *marketplace.Install, db database.Datastore) error {

	if meta.Backend == "terraform" {

		err := terraform.AddApplicationToStack(appConfig, location, meta, install, db)
		app := models.Application{
			Name:    install.Applications.Name,
			Version: install.Applications.Version,
			Source:  meta.Package,
			Status:  "STAGED",
		}
		deployment := models.Deployment{
			Name:         install.DeploymentName,
			Applications: []models.Application{app},
			Creator:      "admin",
			CreationDate: time.Time{},
		}
		deploymentID, err := db.StoreDeployment(deployment)
		if err != nil {
			return err
		}
		if err != nil {
			db.UpdateApplicationStatus(deploymentID, install.Applications.Name, "STAGINGFAILED")

			return err
		}

		executor := &terraform.RealTerraformExecutor{}

		//m, err := fetchMandatoryVars()
		//if err != nil {
		//	return err
		//}
		//terraform.WriteTFVars(m, appConfig)
		err = runPreInstall(appConfig, meta)
		if err != nil {
			return err
		}
		err = terraform.RunTerraform(appConfig, conn, userid, executor)
		if err != nil {
			db.UpdateApplicationStatus(deploymentID, install.Applications.Name, "FAILED")
			return err
		}
		db.UpdateApplicationStatus(deploymentID, install.Applications.Name, "INSTALLED")
		err = runPostInstall(appConfig, meta)

		if err != nil {
			db.UpdateApplicationStatus(deploymentID, install.Applications.Name, "POSTINSTALL FAILED")

			return err
		}
		db.UpdateApplicationStatus(deploymentID, install.Applications.Name, "COMPLETE")

		return nil
	} else {
		return errors.New("backend not implemented")
	}
}

func runPostInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata) error {

	if meta.PostInstall != "" {
		//TODO UNPIN ME
		cmd := exec.Command(filepath.Join(appConfig.Workdir, "workspace", meta.Name, meta.PostInstall))
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("NAME=%s", meta.Name))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", meta.WorkDirectory))
		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running post install script")
			return err
		}
	}
	return nil
}

func runPreInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata) error {
	if meta.PreInstall != "" {
		//TODO UNPIN ME
		cmd := exec.Command(filepath.Join(appConfig.Workdir, "workspace", meta.Name, meta.PostInstall))
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("NAME=%s", meta.Name))
		cmd.Env = append(cmd.Env, fmt.Sprintf("WORKDIR=%s", meta.WorkDirectory))
		cmd.Env = append(cmd.Env, fmt.Sprintf("EKS_NAME=%s", "test_deployment"))
		if err := cmd.Run(); err != nil {
			log.WithError(err).Error("Error running post install script")
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
