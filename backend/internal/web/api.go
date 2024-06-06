package web

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	// "github.com/unity-sds/unity-management-console/backend/internal/aws"
	// "github.com/unity-sds/unity-management-console/backend/internal/database"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
	// strftime "github.com/ncruces/go-strftime"
	"net/http"
	// "time"
)

func handleAPICall(appConfig config.AppConfig) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Get the location of the health check bucket
		// healthCheckParamPath := fmt.Sprintf("/unity/deployment/%s/%s/cs/monitoring/s3/bucketName", appConfig.Project, appConfig.Venue)
		// bucketNameParam, err := aws.ReadSSMParameter(healthCheckParamPath)

		// // Get a listing of all the files in the bucket and pick the one with the latest timestamp
		// result := aws.ListObjectsV2(nil, &appConfig, *bucketNameParam.Parameter.Value, "health_check")

		// layout, err := strftime.Layout("health_check_%Y-%m-%d_%H-%M-%S.json")
		// if err != nil {
		// 	log.Warnf("%s", "Error parsing date layout")
		// }

		// var latestHealthCheckObject *types.Object
		// var latestHealthCheckDatetime *time.Time

		// for i, object := range result {
		// 	t, err := time.Parse(layout, *object.Key)

		// 	if err != nil || t.IsZero() {
		// 		log.Warnf("File Doesn't Match: %s", *object.Key)
		// 		continue
		// 	}

		// 	if latestHealthCheckObject == nil || t.After(*latestHealthCheckDatetime) {
		// 		latestHealthCheckObject = &result[i]
		// 		latestHealthCheckDatetime = &t
		// 	}
		// }

		// if latestHealthCheckObject == nil {
			jsonData := []byte(`{"error": "Can't find any health check files"}`)
			c.Data(http.StatusOK, "application/json", jsonData)
		// }

		// // Read the object and pass the data on to the requester
		// object := aws.GetObject(nil, &appConfig, *bucketNameParam.Parameter.Value, *latestHealthCheckObject.Key)
		// c.Data(http.StatusOK, "application/json", object)

		return
	}
	return gin.HandlerFunc(fn)
}
