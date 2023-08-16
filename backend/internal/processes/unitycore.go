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
	"github.com/zclconf/go-cty/cty"
	"os"
	"os/exec"
	"path/filepath"
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
func InstallMarketplaceApplication(conn *websocket.WebSocketManager, userid string, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, location string, install *marketplace.Install) error {

	if meta.Backend == "terraform" {

		err := terraform.AddApplicationToStack(appConfig, location, meta, install)
		if err != nil {
			return err
		}

		executor := &terraform.RealTerraformExecutor{}

		m, err := fetchMandatoryVars()
		if err != nil {
			return err
		}
		terraform.WriteTFVars(m, appConfig)
		err = runPreInstall(meta)
		if err != nil {
			return err
		}
		err = terraform.RunTerraform(appConfig, conn, userid, executor)
		if err != nil {
			return err
		}

		err = runPostInstall(appConfig, meta)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("backend not implemented")
	}
}

func runPostInstall(appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata) error {

	if meta.PostInstall != "" {
		cmd := exec.Command(filepath.Join(appConfig.Workdir, "terraform", "unity-eks", "0.1", meta.PostInstall))
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

func runPreInstall(meta *marketplace.MarketplaceMetadata) error {

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

	if err := InstallMarketplaceApplication(wsManager, userid, conf, meta, location, received); err != nil {
		log.Error("Error installing application:", err)
		return err
	}

	return nil
}
