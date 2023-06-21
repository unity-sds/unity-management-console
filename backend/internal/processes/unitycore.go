package processes

import (
	"github.com/go-git/go-git/v5"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/act"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"github.com/unity-sds/unity-control-plane/backend/internal/metadata"
	"os"
	"strings"
)

type ActRunner interface {
	RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error
}

type ActRunnerImpl struct{}

func (r *ActRunnerImpl) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error {
	return act.RunAct(path, inputs, env, secrets, conn, appConfig)
}

func (r *ActRunnerImpl) UpdateCoreConfig(conn *websocket.Conn, store database.Datastore, config config.AppConfig) error {
	inputs := map[string]string{
		"deploymentProject": "SIPS",
		"deploymentStage":   "SIPS",
		"awsConnection":     "keys",
	}
	cParams, err := store.FetchCoreParams()
	if err != nil {
		log.Errorf("Error fetching params. %v", err)
		return err
	}
	project := ""
	venue := ""
	privsubs := ""
	pubsubs := ""
	for _, v := range cParams {
		if v.Key == "project" {
			project = v.Value
		} else if v.Key == "venue" {
			venue = v.Value
		} else if v.Key == "privateSubnets" {
			privsubs = v.Value
		} else if v.Key == "publicSubnets" {
			pubsubs = v.Value
		}
	}
	env := map[string]string{
		"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
		"AWS_REGION":            "us-west-2",
		"CORE_PROJECT":          project,
		"CORE_VENUE":            venue,
		"CORE_PRIVATE_SUBNETS":  privsubs,
		"CORE_PUBLIC_SUBNETS":   pubsubs,
	}

	secrets := map[string]string{}
	return r.RunAct(config.WorkflowBasePath+"/environment-provisioner.yml", inputs, env, secrets, conn, config)
}

func (r *ActRunnerImpl) ValidateMarketplaceInstallation(name string, version string) (bool, marketplace.MarketplaceMetadata, error) {
	// Validate installation

	// Check Marketplace Installation Exists
	meta, err := fetchMarketplaceMetadata(name, version)
	if err != nil {
		return false, meta, err
	}

	// Is already installed?
	exists, err := checkExistingInstallation(meta)
	if exists == true || err != nil {
		return false, meta, err
	}

	// Do dependencies match?
	depvalid, err := checkDependencies(meta)
	if depvalid == false || err != nil {
		return false, meta, err
	}

	return true, meta, nil
}

func (r *ActRunnerImpl) FetchPackage(meta marketplace.MarketplaceMetadata) (string, error) {
	// Get package

	locationdir := ""
	if strings.HasSuffix(meta.Package, ".zip") {
		// Fetch from zip

	} else {
		// Checkout git repo
		locationdir, err := gitclone(meta.Package)
		return locationdir, err
	}
	return locationdir, nil
}

func gitclone(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "git-")
	if err != nil {
		return tempDir, err
	}
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return tempDir, err
}

func (r *ActRunnerImpl) GenerateMetadata(appname string, loc string, extensions *marketplace.Install_Extensions) (string, error) {
	// Generate meta string
	if appname == "unity-eks" {
		meta, err := metadata.GenerateEKSMetadata(extensions)
		return string(meta), err
	}
	return "", nil
}

func (r *ActRunnerImpl) CheckIAMPolicies() error {
	// Check IAM policies

	// Get default polcies from marketplace

	// Run IAM Simulator
	return nil
}
func (r *ActRunnerImpl) InstallMarketplaceApplication(conn *websocket.Conn, store database.Datastore, meta string, config config.AppConfig, entrypoint string) error {

	// Install package
	inputs := map[string]string{
		"META":             meta,
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

func (r *ActRunnerImpl) TriggerInstall(conn *websocket.Conn, store database.Datastore, received marketplace.Install, conf config.AppConfig) error {

	log.Infof("Token v: %v", conf.GithubToken)
	t := received.Applications
	log.Info("Validating installation")
	validmarket, meta, err := r.ValidateMarketplaceInstallation(t.Name, t.Version)
	if err != nil || validmarket == false {
		return err
	}
	log.Info("Checking IAM Policies")
	err = r.CheckIAMPolicies()
	if err != nil {
		return err
	}
	log.Info("Fetching package")
	location, err := r.FetchPackage(meta)
	if err != nil {
		return err
	}
	log.Info("Generating Metadata")
	metastr, err := r.GenerateMetadata(t.Name, location, received.Extensions)
	if err != nil {
		return err
	}

	log.Info("Installing Application")
	err = r.InstallMarketplaceApplication(conn, store, metastr, conf, meta.Entrypoint)
	if err != nil {
		return err
	}

	return nil
}
