package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	// "github.com/unity-sds/unity-management-console/backend/internal/database"
	"net/http"
)

func handleAPICall(appConfig config.AppConfig) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get the location of the health check bucket
		healthCheckParamPath := fmt.Sprintf("/unity/deployment/%s/%s/cs/monitoring/s3/bucketName", appConfig.Project, appConfig.Venue)
		bucketNameParam, err := aws.ReadSSMParameter(healthCheckParamPath)

		// Get a listing of all the files in the bucket and pick the one with the latest timestamp
		result := aws.ListObjectsV2(nil, &appConfig, *bucketNameParam.Parameter.Value, "health_check")

		for _, object := range result {
			log.Warnf("%v", *object.Key)	
		}
		
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
