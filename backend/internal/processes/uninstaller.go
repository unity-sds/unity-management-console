package processes

import (
	"bufio"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type UninstallPayload struct {
	Application        string
	ApplicationPackage string
	Deployment         string
}

func UninstallApplication(payload string, conf *config.AppConfig, store database.Datastore) error {

	filepath := path.Join(conf.Workdir, "workspace")
	var uninstall UninstallPayload
	err := json.Unmarshal([]byte(payload), &uninstall)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Checking file	%s has prefix: %s", file.Name(), uninstall.ApplicationPackage)
		if strings.HasPrefix(file.Name(), uninstall.ApplicationPackage) {
			log.Infof("File was a match")
			// Open the file
			f, err := os.Open(path.Join(filepath, file.Name()))
			if err != nil {
				return err
			}

			// Read comments at the top
			scanner := bufio.NewScanner(f)
			metadata := make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				if !strings.HasPrefix(line, "#") {
					break
				}
				// Parsing the comments
				parts := strings.SplitN(strings.TrimPrefix(line, "# "), ": ", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := strings.TrimSpace(parts[1])
					metadata[key] = value
				}
			}
			f.Close()

			// Check applicationName from the comments and delete the file if it matches
			log.Infof("Check if appname %s == %s", metadata["applicationName"], uninstall.Application)
			if metadata["applicationName"] == uninstall.Application {
				err = os.Remove(path.Join(filepath, file.Name()))
				if err != nil {
					id, err := store.FetchDeploymentIDByName(uninstall.Deployment)
					err = store.UpdateApplicationStatus(id, uninstall.Application, "UNINSTALL FAILED")
					return err
				}
				err := store.RemoveApplicationByName(uninstall.Deployment, uninstall.Application)
				if err != nil {
					id, err := store.FetchDeploymentIDByName(uninstall.Deployment)
					err = store.UpdateApplicationStatus(id, uninstall.Application, "UNINSTALL FAILED")
					return err
				}
				err = fetchAllApplications(store)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return nil
}

func ReapplyApplication(payload string, conf *config.AppConfig, store database.Datastore, wsmgr *websocket.WebSocketManager, userid string) error {
	filepath := path.Join(conf.Workdir, "workspace")
	var uninstall UninstallPayload
	err := json.Unmarshal([]byte(payload), &uninstall)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Checking file	%s has prefix: %s", file.Name(), uninstall.ApplicationPackage)
		if strings.HasPrefix(file.Name(), uninstall.ApplicationPackage) {
			log.Infof("File was a match")
			// Open the file
			f, err := os.Open(path.Join(filepath, file.Name()))
			if err != nil {
				return err
			}

			// Read comments at the top
			scanner := bufio.NewScanner(f)
			metadata := make(map[string]string)
			for scanner.Scan() {
				line := scanner.Text()
				if !strings.HasPrefix(line, "#") {
					break
				}
				// Parsing the comments
				parts := strings.SplitN(strings.TrimPrefix(line, "# "), ": ", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := strings.TrimSpace(parts[1])
					metadata[key] = value
				}
			}
			f.Close()

			// Check applicationName from the comments and delete the file if it matches
			log.Infof("Check if appname %s == %s", metadata["applicationName"], uninstall.Application)
			if metadata["applicationName"] == uninstall.Application {
				app := marketplace.Install_Applications{
					Name:        metadata["application"],
					Version:     metadata["version"],
					Variables:   nil,
					Postinstall: "",
					Preinstall:  "",
				}
				meta, err := validateAndPrepareInstallation(&app, conf)
				if err != nil {
					return err
				}
				val, err := strconv.ParseUint(metadata["deploymentID"], 10, 0)
				uintVal := uint(val)
				execute(store, conf, meta, metadata["applicationName"], uintVal, wsmgr, userid, metadata["deploymentID"])
			}
		}
	}

	return nil
}
