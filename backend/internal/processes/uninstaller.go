package processes

import (
	"bufio"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"fmt"
)

type UninstallPayload struct {
	Application        string
	DisplayName        string
	ApplicationPackage string
	Deployment         string
}

func UninstallAll(conf *config.AppConfig, conn *websocket.WebSocketManager, userid string, received *marketplace.Uninstall) error {
	executor := &terraform.RealTerraformExecutor{}
	err := terraform.DestroyAllTerraform(conf, conn, userid, executor)
	if err != nil {
		log.WithError(err).Error("FAILED TO DESTROY ALL COMPONENTS")
		//return err
	}

	if (received.DeleteBucket) {
		err = aws.DeleteS3Bucket(conf.BucketName)
		if err != nil {
			log.WithError(err).Error("FAILED TO REMOVE S3 BUCKET")
		}		
	}

	err = aws.DeleteStateTable(conf.InstallPrefix)
	if err != nil {
		log.WithError(err).Error("FAILED TO REMOVE DYNAMODB TABLE")
	}

	log.Info("UNITY MANAGEMENT CONSOLE UNINSTALL COMPLETE")
	return nil
}

func UninstallApplication(appname string, deploymentname string, displayname string, conf *config.AppConfig, store database.Datastore, conn *websocket.WebSocketManager, userid string) error {
	executor := &terraform.RealTerraformExecutor{}

	filepath := path.Join(conf.Workdir, "workspace")

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Checking file %s has prefix: %s", file.Name(), appname)
		if strings.HasPrefix(file.Name(), appname) {
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
			log.Infof("Check if appname %s == %s", metadata["applicationName"], displayname)
			if metadata["applicationName"] == displayname {
				p := path.Join(filepath, file.Name())
				log.Infof("Attempting to delete file: %s", p)
				err = os.Remove(p)
				if err != nil {
					id, err := store.FetchDeploymentIDByName(deploymentname)
					log.WithError(err).Error("Failed to fetch deployment ID by name when removing application")
					err = store.UpdateApplicationStatus(id, appname, displayname, "UNINSTALL FAILED")
					log.WithError(err).Error("Failed to update application status removing application")
					return err
				}
				err := store.RemoveApplicationByName(deploymentname, appname)
				if err != nil {
					id, err := store.FetchDeploymentIDByName(deploymentname)
					log.WithError(err).Error("Failed to fetch deployment ID by name when removing application")
					err = store.UpdateApplicationStatus(id, appname, displayname, "UNINSTALL FAILED")
					log.WithError(err).Error("Failed to update application status removing application")
					return err
				}
				err = fetchAllApplications(store)
				if err != nil {
					return err
				}
				err = terraform.RunTerraform(conf, conn, userid, executor, "")
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return nil
}

func UninstallApplicationNew(appname string, deploymentname string, displayname string, conf *config.AppConfig, store database.Datastore) error {
	// Create uninstall_logs directory if it doesn't exist
	logDir := path.Join(conf.Workdir, "uninstall_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create install_logs directory: %w", err)
	}

	executor := &terraform.RealTerraformExecutor{}

	filepath := path.Join(conf.Workdir, "workspace")

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Checking file %s has prefix: %s", file.Name(), appname)
		if strings.HasPrefix(file.Name(), appname) {
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
			log.Infof("Check if appname %s == %s", metadata["applicationName"], displayname)
			if metadata["applicationName"] == displayname {
				p := path.Join(filepath, file.Name())
				log.Infof("Attempting to delete file: %s", p)
				err = os.Remove(p)
				if err != nil {
					id, err := store.FetchDeploymentIDByName(deploymentname)
					log.WithError(err).Error("Failed to fetch deployment ID by name when removing application")
					err = store.UpdateApplicationStatus(id, appname, displayname, "UNINSTALL FAILED")
					log.WithError(err).Error("Failed to update application status removing application")
					return err
				}
				logfile := path.Join(logDir, fmt.Sprintf("%s_uninstall_log", deploymentname))
				err = terraform.RunTerraformLogOutToFile(conf, logfile, executor, "")
				if err != nil {
					log.WithError(err).Error("Failed to uninstall application")
					return err
				}

				// err := store.RemoveApplicationByName(deploymentname, appname)
				// if err != nil {
				// 	id, err := store.FetchDeploymentIDByName(deploymentname)
				// 	log.WithError(err).Error("Failed to fetch deployment ID by name when removing application")
				// 	err = store.UpdateApplicationStatus(id, appname, displayname, "UNINSTALL FAILED")
				// 	log.WithError(err).Error("Failed to update application status removing application")
				// 	return err
				// }
				id, err := store.FetchDeploymentIDByName(deploymentname)
				err = store.UpdateApplicationStatus(id, appname, displayname, "UNINSTALLED")
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

func UninstallApplicationNewV2(appName string, version string, deploymentName string, conf *config.AppConfig, store database.Datastore) error {
	// Create uninstall_logs directory if it doesn't exist
	logDir := path.Join(conf.Workdir, "uninstall_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create install_logs directory: %w", err)
	}

	executor := &terraform.RealTerraformExecutor{}

	filepath := path.Join(conf.Workdir, "workspace")

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Checking file %s has prefix: %s", file.Name(), appName)
		if strings.HasPrefix(file.Name(), appName) {
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
			log.Infof("Check if appname %s == %s", metadata["applicationName"], deploymentName)
			if metadata["applicationName"] == deploymentName {
				p := path.Join(filepath, file.Name())
				log.Infof("Attempting to delete file: %s", p)
				err = os.Remove(p)
				// if err != nil {
				// 	log.WithError(err).Error("Failed to fetch deployment ID by name when removing application")
				// 	err = store.UpdateApplicationStatus(id, appname, deploymentName, "UNINSTALL FAILED")
				// 	log.WithError(err).Error("Failed to update application status removing application")
				// 	return err
				// }
				store.UpdateInstalledMarketplaceApplicationStatusByName(appName, deploymentName, "STARTING UNINSTALL")
				logfile := path.Join(logDir, fmt.Sprintf("%s_%s_uninstall_log", appName, deploymentName))
				err = terraform.RunTerraformLogOutToFile(conf, logfile, executor, "")
				if err != nil {
					log.WithError(err).Error("Failed to uninstall application")
					return err
				}

				err = store.RemoveInstalledMarketplaceApplicationByName(appName)

				// err := store.RemoveApplicationByName(deploymentname, appname)
				// if err != nil {
				// 	id, err := store.FetchDeploymentIDByName(deploymentname)
				// 	log.WithError(err).Error("Failed to fetch deployment ID by name when removing application")
				// 	err = store.UpdateApplicationStatus(id, appname, deploymentName, "UNINSTALL FAILED")
				// 	log.WithError(err).Error("Failed to update application status removing application")
				// 	return err
				// }
				// id, err := store.FetchDeploymentIDByName(deploymentname)
				err = store.UpdateInstalledMarketplaceApplicationStatusByName(appName, deploymentName, "UNINSTALLED")
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
				inst := &marketplace.Install{
					Applications:   nil,
					DeploymentName: metadata["deploymentID"],
				}
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
				err = execute(store, conf, meta, inst, uintVal, wsmgr, userid)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
