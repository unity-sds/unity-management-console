package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"net/http"
)

func handleAPICall(appConfig config.AppConfig) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		healthCheckParamPath := fmt.Sprintf("/unity/%s/%s/cs/monitoring/s3/bucketName", appConfig.Project, appConfig.Venue)
		param, err := aws.ReadSSMParameter(healthCheckParamPath)

		if err != nil {
			log.WithError(err).Error("Failed to get SSM params for " + healthCheckParamPath)
			c.JSON(http.StatusInternalServerError, "")
		}
		fmt.Printf("%v", param)

		c.JSON(http.StatusOK, gin.H{
			"health_checks": "not done",
		})

	}
	return gin.HandlerFunc(fn)
}
