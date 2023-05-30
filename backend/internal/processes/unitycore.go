package processes

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/act"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"os"
)

var basepath = "/home/barber/Projects/unity-cs-infra/"

type ActRunner interface {
	RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error
}

type ActRunnerImpl struct{}

func (r *ActRunnerImpl) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error {
	return act.RunAct(path, inputs, env, secrets, conn)
}
func fetchCoreParams() ([]models.CoreConfig, error) {
	config, err := database.FetchConfig()
	if err != nil {
		return []models.CoreConfig{}, nil
	} else {
		return config, nil
	}
}
func UpdateCoreConfig(conn *websocket.Conn, store database.Datastore, runner ActRunner) error {
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
	//act.RunAct(basepath+".github/workflows/environment-provisioner.yml", inputs, env, secrets, conn)
	return runner.RunAct(basepath+".github/workflows/environment-provisioner.yml", inputs, env, secrets, conn)
}
