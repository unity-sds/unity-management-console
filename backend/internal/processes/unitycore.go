package processes

import (
	"encoding/base64"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/act"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"github.com/unity-sds/unity-control-plane/backend/internal/metadata"
	"os"
)

type ActRunner interface {
	RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error
}

type ActRunnerImpl struct {
}

// NewActRunner creates a new ActRunnerImpl instance.
func NewActRunner() *ActRunnerImpl {
	return &ActRunnerImpl{}
}
func (r *ActRunnerImpl) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error {
	return act.RunAct(path, inputs, env, secrets, conn, appConfig)
}

func GenerateMetadata(appname string, install *marketplace.Install, meta *marketplace.MarketplaceMetadata) ([]byte, error) {
	// Generate meta string
	if appname == "unity-eks" {
		meta, err := metadata.GenerateEKSMetadata(install.Extensions)
		return meta, err
	}

	metaarr, err := metadata.GenerateApplicationMetadata(appname, install, meta)
	return metaarr, err
}

func InstallMarketplaceApplication(conn *websocket.Conn, meta []byte, config config.AppConfig, entrypoint string, r ActRunnerImpl, location string) error {

	str := base64.StdEncoding.EncodeToString(meta)
	// Install package
	inputs := map[string]string{
		"META":             str,
		"DEPLOYMENTSOURCE": "act",
		"AWSCONNECTION":    "keys",
	}

	env := map[string]string{
		"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
		"AWS_REGION":            "us-west-2",
	}

	secrets := map[string]string{
		"token": config.GithubToken,
	}
	log.Infof("Launching act runner with following meta: %v", meta)
	action := config.WorkflowBasePath + "/install-stacks.yml"
	if entrypoint != "" {
		action = config.WorkflowBasePath + "/" + entrypoint
	}

	return r.RunAct(action, inputs, env, secrets, conn, config)

	// Add application to installed packages in database

}

func TriggerInstall(conn *websocket.Conn, store database.Datastore, received marketplace.Install, conf config.AppConfig, r ActRunnerImpl) error {
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

	metastr, err := GenerateMetadata(t.Name, &received, meta)
	if err != nil {
		log.Error("Error generating metadata:", err)
		return err
	}

	if err := InstallMarketplaceApplication(conn, metastr, conf, meta.Entrypoint, r, location); err != nil {
		log.Error("Error installing application:", err)
		return err
	}

	return nil
}
