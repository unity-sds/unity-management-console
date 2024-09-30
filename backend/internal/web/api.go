package web

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	"net/http"
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

func handleApplicationInstall(appConfig config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {

		type ApplicationInstallParams struct {
			Name      string
			Version   string
			DeploymentName string
			Variables string
		}

		var applicationInstallParams ApplicationInstallParams
		err := c.BindJSON(&applicationInstallParams)

		if err != nil {
			log.Errorf("Error parsing JSON: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Bad JSON"})
			return
		}

		log.Errorf("Got JSON: %v", applicationInstallParams)
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
