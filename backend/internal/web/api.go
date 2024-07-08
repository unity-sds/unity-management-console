package web

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	strftime "github.com/ncruces/go-strftime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"net/http"
	"time"
	// "fmt"
)

func handleHealthChecks(c *gin.Context, appConfig config.AppConfig) {
	bucketname := viper.Get("bucketname").(string)

	// Get a listing of all the files in the bucket and pick the one with the latest timestamp
	result := aws.ListObjectsV2(nil, &appConfig, bucketname, "health_check")

	layout, err := strftime.Layout("health_check_%Y-%m-%d_%H-%M-%S.json")
	if err != nil {
		log.Warnf("%s", "Error parsing date layout")
	}

	var latestHealthCheckObject *types.Object
	var latestHealthCheckDatetime *time.Time

	for i, object := range result {
		t, err := time.Parse(layout, *object.Key)

		if err != nil || t.IsZero() {
			log.Warnf("File Doesn't Match: %s", *object.Key)
			continue
		}

		if latestHealthCheckObject == nil || t.After(*latestHealthCheckDatetime) {
			latestHealthCheckObject = &result[i]
			latestHealthCheckDatetime = &t
		}
	}

	if latestHealthCheckObject == nil {
		jsonData := []byte(`{"error": "Can't find any health check files"}`)
		c.Data(http.StatusOK, "application/json", jsonData)
		return
	}

	// Read the object and pass the data on to the requester
	object := aws.GetObject(nil, &appConfig, bucketname, *latestHealthCheckObject.Key)
	c.Data(http.StatusOK, "application/json", object)
}

func handleUninstall(c *gin.Context, appConfig config.AppConfig) {
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

func handleGetAPICall(appConfig config.AppConfig) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		switch endpoint := c.Param("endpoint"); endpoint {
		case "health_checks":
			handleHealthChecks(c, appConfig)
		default:
			handleNoRoute(c)
		}
	}
	return gin.HandlerFunc(fn)
}

func handlePostAPICall(appConfig config.AppConfig) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		switch endpoint := c.Param("endpoint"); endpoint {
		case "uninstall":
			handleUninstall(c, appConfig)
		default:
			handleNoRoute(c)
		}
	}
	return gin.HandlerFunc(fn)
}
