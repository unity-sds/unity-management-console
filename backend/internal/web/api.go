package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	"github.com/unity-sds/unity-management-console/backend/internal/update"
	"github.com/unity-sds/unity-management-console/backend/types"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

		// Kick off install process in async mode. Errors will come back from the initial checks, otherwise we can return OK to user.
		err = processes.TriggerInstall(db, &applicationInstallParams, &appConfig, false)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

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
		deploymentName := c.Param("deploymentName")

		app, err := db.GetInstalledMarketplaceApplication(appName, deploymentName)

		if err != nil {
			log.Errorf("Installed application not found: %v", err)
			c.Status(http.StatusNotFound)
			return
		}

		go processes.UninstallApplication(app, &conf, db)
	}
}

func handleGetApplicationInstallStatusByName(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		appName := c.Param("appName")
		deploymentName := c.Param("deploymentName")
		app, err := db.GetInstalledMarketplaceApplication(appName, deploymentName)

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

func handleDeleteApplication(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		appName := c.Param("appName")
		deploymentName := c.Param("deploymentName")

		existingApplication, err := db.GetInstalledMarketplaceApplication(appName, deploymentName)
		if existingApplication == nil {
			log.Errorf("Unable to find application %s and deployment %s", appName, deploymentName)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Application or deployment name doesn't exist."})
			return
		}

		err = db.RemoveInstalledMarketplaceApplication(appName, deploymentName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	}
}

func handleConfigRequest(appConfig config.AppConfig, db database.Datastore) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get the public and private subnets
		pub, priv, err := aws.FetchSubnets()
		if err != nil {
			log.WithError(err).Error("Error fetching subnets")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subnet information"})
			return
		}

		// Get audit information
		audit, _ := db.FindLastAuditLineByOperation(application.Config_Updated)
		bootstrapFailed, _ := db.FindLastAuditLineByOperation(application.Bootstrap_Unsuccessful)
		bootstrapSuccess, _ := db.FindLastAuditLineByOperation(application.Bootstrap_Successful)

		// Determine bootstrap status
		bootstrapStatus := ""
		if bootstrapSuccess.Owner != "" {
			bootstrapStatus = "complete"
		} else if bootstrapFailed.Owner != "" {
			bootstrapStatus = "failed"
		}

		// Create response object
		configResponse := gin.H{
			"applicationConfig": gin.H{
				"MarketplaceOwner": appConfig.MarketplaceOwner,
				"MarketplaceUser":  appConfig.MarketplaceRepo,
				"Project":           appConfig.Project,
				"Venue":             appConfig.Venue,
				"Version":           appConfig.Version,
			},
			"networkConfig": gin.H{
				"publicsubnets":  pub,
				"privatesubnets": priv,
			},
			"lastupdated": audit.CreatedAt.Format("2006-01-02T15:04:05.000"),
			"updatedby":   audit.Owner,
			"bootstrap":    bootstrapStatus,
			"version":      appConfig.Version,
		}

		c.JSON(http.StatusOK, configResponse)
	}
}

// handleUpdateManagementConsole downloads the latest release and copies all files to the current directory
func handleUpdateManagementConsole(appConfig config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Info("Starting management console update process")

		err := update.UpdateManagementConsoleInPlace()
		if err != nil {
			log.WithError(err).Error("Failed to install release: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"error":   nil,
		})
		return
	}
}

func handleCheckAppDependencies(appConfig config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		appName := c.Param("appName")
		version := c.Param("version")

		metadata, err := processes.FetchMarketplaceMetadata(appName, version, &appConfig)
		if err != nil {
			log.Errorf("Unable to fetch metadata for application: %s, %v", appName, err)
			log.WithError(err).Error("Unable to fetch package")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		errors := false
		results := make(map[string]string)
		for label, ssmParam := range metadata.GetDependencies() {
			formattedParam := strings.Replace(ssmParam, "${PROJ}", appConfig.Project, -1)
			formattedParam = strings.Replace(formattedParam, "${VENUE}", appConfig.Venue, -1)
			param, err := aws.ReadSSMParameter(formattedParam)

			if err != nil {
				log.WithError(err).Error("Unable to get SSM param.")
				results[label] = ""
				errors = true
				continue
			}
			results[label] = *param.Parameter.Value
		}

		if errors {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Missing SSM Parameters",
				"params":  results,
			})
			return
		}

		log.Info("Checking dependencies for %s, version %s", appName, version)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"error":   nil,
			"params":  results,
		})
		return
	}
}
