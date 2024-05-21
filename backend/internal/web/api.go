package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	// "github.com/unity-sds/unity-management-console/backend/internal/database"
	"net/http"
	"time"
	strftime "github.com/ncruces/go-strftime"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func handleAPICall(appConfig config.AppConfig) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get the location of the health check bucket
		healthCheckParamPath := fmt.Sprintf("/unity/deployment/%s/%s/cs/monitoring/s3/bucketName", appConfig.Project, appConfig.Venue)
		bucketNameParam, err := aws.ReadSSMParameter(healthCheckParamPath)

		// Get a listing of all the files in the bucket and pick the one with the latest timestamp
		result := aws.ListObjectsV2(nil, &appConfig, *bucketNameParam.Parameter.Value, "health_check")
		
		layout, err := strftime.Layout("health_check_%Y-%m-%d_%H-%M-%S.json")
		if err != nil {
			log.Warnf("%s", "Error parsing date layout")					
		}

		var latestHealthCheckObject *types.Object
		var latestHealthCheckDatetime *time.Time

		for _, object := range result {
				t, err := time.Parse(layout, *object.Key)
				
				if err != nil  {
					// log.Warnf("File Doesn't Match: %s", *object.Key)
					continue
				}

				if latestHealthCheckObject == nil || t.After(*latestHealthCheckDatetime) {
					latestHealthCheckObject = &object
					latestHealthCheckDatetime = &t
				}
		}

		if latestHealthCheckObject == nil {
			jsonData := []byte(`{"error": "Can't find any health heck files"}`)
			c.Data(http.StatusOK, "application/json", jsonData)
		}



		// Read the object and pass the data on to the requester
		log.Warnf("%v", *latestHealthCheckObject.Key)
		object := aws.GetObject(nil, &appConfig, *bucketNameParam.Parameter.Value, *latestHealthCheckObject.Key)
		log.Warnf("%v", object.Body)
		
		return



		// db, err := database.NewGormDatastore()

		// params, err := db.FetchSSMParams()

		// p, err := aws.ReadSSMParameters(params)

		if err != nil {
			log.WithError(err).Error("Failed to get SSM params for " + healthCheckParamPath)
		}
		// fmt.Printf("%v", bucketName)

		aws.GetObject(nil, &appConfig, *bucketNameParam.Parameter.Value, "")

		outStr := fmt.Sprintf(`{"msg":"%s"}`, *bucketNameParam.Parameter.Value)
		jsonData := []byte(outStr)
		c.Data(http.StatusOK, "application/json", jsonData)

	}
	return gin.HandlerFunc(fn)
}
