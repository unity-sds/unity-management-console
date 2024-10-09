package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	"github.com/unity-sds/unity-management-console/backend/types"
	"net/http"
	"os"
	"path/filepath"
	// "strconv"
)

func handleHealthChecks(appConfig config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		bucketname := viper.Get("bucketname").(string)

		// Get the latest health check file
		result, err := aws.GetObject(nil, &appConfig, bucketname, "health_check_latest.json")
		if err != nil {
			log.Errorf("Error getting health check file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve health check data"})
			return
		}
		c.Data(http.StatusOK, "application/json", result)
	}
}

func handleUninstall(appConfig config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		uninstallStatus := viper.Get("uninstallStatus")

		if uninstallStatus != nil {
			c.JSON(http.StatusOK, gin.H{"uninstall_status": uninstallStatus})
			return
		}

		var uninstallOptions struct {
			DeleteBucket *bool `form:"delete_bucket" json:"delete_bucket"`
		}
		err := c.BindJSON(&uninstallOptions)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad input posted."})
			return
		}

		deleteBucket := false
		if uninstallOptions.DeleteBucket != nil {
			deleteBucket = *uninstallOptions.DeleteBucket
		}

		received := &marketplace.Uninstall{
			DeleteBucket: deleteBucket,
		}

		go processes.UninstallAll(&conf, nil, "restAPIUser", received)
		viper.Set("uninstallStatus", "in progress")
		c.JSON(http.StatusOK, gin.H{"uninstall_status": "in progress"})
	}
}

func handleApplicationInstall(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		var applicationInstallParams types.ApplicationInstallParams
		err := c.BindJSON(&applicationInstallParams)

		if err != nil {
			log.Errorf("Error parsing JSON: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Bad JSON"})
			return
		}

		log.Errorf("Got JSON: %v", applicationInstallParams)

		// First check if this application is already installed.
		existingApplication, err := db.GetInstalledMarketplaceApplicationStatusByName(applicationInstallParams.Name, applicationInstallParams.DeploymentName)
		if err != nil {
			log.WithError(err).Error("Error finding applications")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to search applcation list"})	
		}
		
		if existingApplication != nil && existingApplication.Status != "UNINSTALLED" {
			errMsg := fmt.Sprintf("Application with name %s already exists. Status: %s", applicationInstallParams.Name, existingApplication.Status)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})	
			return
		}

		// OK to start installing, get the metadata for this application
		metadata, err := processes.FetchMarketplaceMetadata(applicationInstallParams.Name, applicationInstallParams.Version, &appConfig)
		if err != nil {
			log.Errorf("Unable to fetch metadata for application: %s, %v", applicationInstallParams.Name, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch package"})
			return
		}

		// Install the application package files
		location, err := processes.FetchPackage(&metadata, &appConfig)
		if err != nil {
			log.Errorf("Unable to fetch package for application: %s, %v", applicationInstallParams.Name, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch package"})
			return
		}

		err = processes.InstallMarketplaceApplicationNewV2(&appConfig, location, &applicationInstallParams, &metadata, db, false)
		
		c.Status(http.StatusOK)

	}
}

func handleGetInstallLogs(appConfig config.AppConfig, db database.Datastore, uninstall bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		appName := c.Param("appName")
		deploymentName := c.Param("deploymentName")

		// deploymentID, err := db.FetchDeploymentIDByApplicationName(deploymentName)
		// if err != nil {
		// 	log.Errorf("Error getting deployment ID: %v", err)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading application status"})
		// 	return
		// }

		var logDir string
		if uninstall {
			logDir = filepath.Join(appConfig.Workdir, "uninstall_logs")
		} else {
			logDir = filepath.Join(appConfig.Workdir, "install_logs")
		}

		var logfile string
		if uninstall {
			logfile = filepath.Join(logDir, fmt.Sprintf("%s_%s_uninstall_log", appName, deploymentName))
		} else {
			logfile = filepath.Join(logDir, fmt.Sprintf("%s_%s_install_log", appName, deploymentName))
		}

		// Read the log file
		content, err := os.ReadFile(logfile)
		if err != nil {
			log.Errorf("Error reading log file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file"})
			return
		}

		// Return the file contents
		c.Data(http.StatusOK, "text/plain", content)
	}
}

func handleUninstallApplication(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		appName := c.Param("appName")
		version := c.Param("version")
		deploymentName := c.Param("deploymentName")

		go processes.UninstallApplicationNewV2(appName, version, deploymentName, &conf, db)
	}
}

func handleGetApplicationInstallStatusByName(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		appName := c.Param("appName")
		deploymentName := c.Param("deploymentName")
		app, err := db.GetInstalledMarketplaceApplicationStatusByName(appName, deploymentName)

		if err != nil {
			log.Errorf("Error reading application status: %v", err)
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, app)
	}
}

func getInstalledApplications(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		applications, err := db.FetchAllInstalledMarketplaceApplications()
		if err != nil {
			log.Errorf("Error getting application list: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, applications)
	}
}

// func handleGetAPICall(appConfig config.AppConfig) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		switch endpoint := c.Param("endpoint"); endpoint {
// 		case "health_checks":
// 			handleHealthChecks(c, appConfig)
// 		default:
// 			handleNoRoute(c)
// 		}
// 	}
// 	return gin.HandlerFunc(fn)
// }

// func handlePostAPICall(appConfig config.AppConfig) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		switch endpoint := c.Param("endpoint"); endpoint {
// 		case "uninstall":
// 			handleUninstall(c, appConfig)
// 		default:
// 			handleNoRoute(c)
// 		}
// 	}
// 	return gin.HandlerFunc(fn)
// }
