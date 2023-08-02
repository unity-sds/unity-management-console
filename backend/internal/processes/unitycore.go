package processes

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
)

func InstallMarketplaceApplication(conn *websocket.WebSocketManager, userid string, appConfig *config.AppConfig, meta *marketplace.MarketplaceMetadata, location string, install *marketplace.Install) error {

	if meta.Backend == "terraform" {

		err := terraform.AddApplicationToStack(appConfig, location, meta)
		if err != nil {
			return err
		}

		executor := &terraform.RealTerraformExecutor{}

		terraform.RunTerraform(appConfig, conn, userid, executor)
		return nil
	} else {
		return errors.New("backend not implemented")
	}
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
