package processes

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"os"
	"strings"
)

func UpdateCoreConfig(conn *websocket.Conn, store database.Datastore, config config.AppConfig, r ActRunner) error {
	inputs := map[string]string{
		"deploymentProject": "SIPS",
		"deploymentStage":   "SIPS",
		"awsConnection":     "keys",
	}
	cParams, err := store.FetchSSMParams()
	if err != nil {
		log.Errorf("Error fetching params. %v", err)
		return err
	}

	env := map[string]string{
		"AWS_ACCESS_KEY_ID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_SESSION_TOKEN":     os.Getenv("AWS_SESSION_TOKEN"),
		"AWS_REGION":            "us-west-2",
	}

	varstr := ""
	for _, v := range cParams {
		if v.Key == "/unity/core/project" {
			env["TF_VAR_project"] = v.Value
		} else if v.Key == "/unity/core/venue" {
			env["TF_VAR_venue"] = v.Value

		} else {
			varstr = varstr + fmt.Sprintf("{name  = \"%v\", type  = \"%v\", value = \"%v\"},", v.Key, v.Type, v.Value)
		}
	}

	if varstr != "" {
		varstr = strings.TrimRight(varstr, ",")
		varstr = fmt.Sprintf("[%v]", varstr)

		env["TF_VAR_ssm_parameters"] = varstr

	}
	secrets := map[string]string{}
	return r.RunAct(config.WorkflowBasePath+"/environment-provisioner.yml", inputs, env, secrets, conn, config)
}
