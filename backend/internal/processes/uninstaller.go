package processes

import (
	"bufio"
	// "encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"github.com/unity-sds/unity-management-console/backend/types"
	"io/ioutil"
	"os"
	"path"
	// "strconv"
	"fmt"
	"os/exec"
	"strings"
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

	if received.DeleteBucket {
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

func runScript(application *types.InstalledMarketplaceApplication, store database.Datastore, path string) error {
	if _, err := os.Stat(path); err == nil {
		application.Status = fmt.Sprintf("RUNNING SCRIPT: %s", path)
		store.UpdateInstalledMarketplaceApplication(application)
		log.Infof("Found script at %s, executing...", path)
		cmd := exec.Command("/bin/sh", path)
		cmd.Env = os.Environ() // Inherit parent environment
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.WithError(err).Errorf("Pre-uninstall script failed: %s", string(output))
			return fmt.Errorf("pre-uninstall script failed: %w", err)
		}
		log.Infof("Pre-uninstall script output: %s", string(output))
	}
	return nil
}

func UninstallApplication(application *types.InstalledMarketplaceApplication, conf *config.AppConfig, store database.Datastore) error {
	// Create uninstall_logs directory if it doesn't exist
	logDir := path.Join(conf.Workdir, "uninstall_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create install_logs directory: %w", err)
	}

	logfile := path.Join(logDir, fmt.Sprintf("%s_%s_uninstall_log", application.Name, application.DeploymentName))
	executor := &terraform.RealTerraformExecutor{}

	application.Status = "UNINSTALLING"
	store.UpdateInstalledMarketplaceApplication(application)

	// Check for and run pre-uninstall script if it exists
	preUninstallScript := path.Join(conf.Workdir, "workspace", application.Name, "pre_uninstall.sh")
	err := runScript(application, store, preUninstallScript)
	if err != nil {
		return err
	}

	// Run a terraform destroy on the module to be uninstalled
	err = terraform.DestroyTerraformModule(conf, logfile, executor, application.TerraformModuleName)
	if err != nil {
		return err
	}

	// Find and delete the module files in our TF workspace
	filepath := path.Join(conf.Workdir, "workspace")

	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Checking file %s has prefix: %s", file.Name(), application.Name)
		if strings.HasPrefix(file.Name(), application.Name) {
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
			log.Infof("Check if appname %s == %s", metadata["applicationName"], application.DeploymentName)
			if metadata["applicationName"] == application.DeploymentName {
				p := path.Join(filepath, file.Name())
				log.Infof("Attempting to delete file: %s", p)
				err = os.Remove(p)

				// Run terraform apply on modified workspace
				logfile := path.Join(logDir, fmt.Sprintf("%s_%s_uninstall_log", application.Name, application.DeploymentName))
				err = terraform.RunTerraformLogOutToFile(conf, logfile, executor, "")
				if err != nil {
					log.WithError(err).Error("Failed to uninstall application")
					return err
				}

				application.Status = "UNINSTALLED"
				store.UpdateInstalledMarketplaceApplication(application)

				err = fetchAllApplications(store)
				if err != nil {
					return err
				}

				// Check for and run pre-uninstall script if it exists
				postUninstallScript := path.Join(conf.Workdir, "workspace", application.Name, "post_uninstall.sh")
				err := runScript(application, store, postUninstallScript)
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
	// filepath := path.Join(conf.Workdir, "workspace")
	// var uninstall UninstallPayload
	// err := json.Unmarshal([]byte(payload), &uninstall)
	// if err != nil {
	// 	return err
	// }

	// files, err := ioutil.ReadDir(filepath)
	// if err != nil {
	// 	return err
	// }

	// for _, file := range files {
	// 	log.Infof("Checking file	%s has prefix: %s", file.Name(), uninstall.ApplicationPackage)
	// 	if strings.HasPrefix(file.Name(), uninstall.ApplicationPackage) {
	// 		log.Infof("File was a match")
	// 		// Open the file
	// 		f, err := os.Open(path.Join(filepath, file.Name()))
	// 		if err != nil {
	// 			return err
	// 		}

	// 		// Read comments at the top
	// 		scanner := bufio.NewScanner(f)
	// 		metadata := make(map[string]string)
	// 		for scanner.Scan() {
	// 			line := scanner.Text()
	// 			if !strings.HasPrefix(line, "#") {
	// 				break
	// 			}
	// 			// Parsing the comments
	// 			parts := strings.SplitN(strings.TrimPrefix(line, "# "), ": ", 2)
	// 			if len(parts) == 2 {
	// 				key := parts[0]
	// 				value := strings.TrimSpace(parts[1])
	// 				metadata[key] = value
	// 			}
	// 		}
	// 		f.Close()

	// 		// Check applicationName from the comments and delete the file if it matches
	// 		log.Infof("Check if appname %s == %s", metadata["applicationName"], uninstall.Application)
	// 		if metadata["applicationName"] == uninstall.Application {
	// 			inst := &marketplace.Install{
	// 				Applications:   nil,
	// 				DeploymentName: metadata["deploymentID"],
	// 			}
	// 			app := marketplace.Install_Applications{
	// 				Name:        metadata["application"],
	// 				Version:     metadata["version"],
	// 				Variables:   nil,
	// 				Postinstall: "",
	// 				Preinstall:  "",
	// 			}
	// 			meta, err := validateAndPrepareInstallation(&app, conf)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			val, err := strconv.ParseUint(metadata["deploymentID"], 10, 0)
	// 			uintVal := uint(val)
	// 			// err = execute(store, conf, meta, inst, uintVal, wsmgr, userid)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

	return nil
}
