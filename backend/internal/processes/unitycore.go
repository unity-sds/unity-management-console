package processes

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/action"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/metadata"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
)

func GenerateMetadata(appname string, install *marketplace.Install, meta *marketplace.MarketplaceMetadata) ([]byte, error) {
	// Generate meta string
	if appname == "unity-eks" {
		return []byte{}, nil
	}

	if appname == "unity-apigateway" {
		return []byte{}, nil
	}
	metaarr, err := metadata.GenerateApplicationMetadata(appname, install, meta)
	return metaarr, err
}

func InstallMarketplaceApplication(conn *websocket.WebSocketManager, userid string, meta []byte, config config.AppConfig, entrypoint string, appName string, install *marketplace.Install) error {

	if appName == "unity-apigateway" {
		return action.RunInstall(conn, userid, install, config)

	} else {
		//str := base64.StdEncoding.EncodeToString(meta)
		// Install package
		//inputs := map[string]string{
		//	"METADATA":         str,
		//	"DEPLOYMENTSOURCE": "act",
		//	"AWSCONNECTION":    "keys",
		//}
		//
		////TODO Figure out how to use packaged workflows from within act runner
		//env := map[string]string{
		//	"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		//	"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		//	"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
		//	"AWS_REGION":            "us-west-2",
		//	"WORKFLOWPATH":          "/home/barber/Projects/unity-management-console/backend/cmd/web/.github/workflows",
		//}
		//
		//secrets := map[string]string{
		//	"token": config.GithubToken,
		//}
		//log.Infof("Launching act runner with following meta: %v", meta)
		//action := config.WorkflowBasePath + "/install-stacks.yml"
		//if entrypoint != "" {
		//	action = config.WorkflowBasePath + "/" + entrypoint
		//}

		//return r.RunAct(action, inputs, env, secrets, conn, config)

		return nil
		// Add application to installed packages in database
	}
}

func TriggerInstall(wsManager *websocket.WebSocketManager, userid string, store database.Datastore, received *marketplace.Install, conf config.AppConfig) error {
	t := received.Applications

	meta, err := validateAndPrepareInstallation(t)
	if err != nil {
		return err
	}

	location, err := FetchPackage(meta)
	if err != nil {
		log.Error("Error fetching package:", err)
		return err
	}

	metastr, err := GenerateMetadata(t.Name, received, meta)
	if err != nil {
		log.Error("Error generating metadata:", err)
		return err
	}

	if err := InstallMarketplaceApplication(wsManager, userid, metastr, conf, meta.Entrypoint, location, received); err != nil {
		log.Error("Error installing application:", err)
		return err
	}

	return nil
}
