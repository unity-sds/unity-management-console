package web

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there was an error
		if len(c.Errors) > 0 {
			// Collect all errors
			var multiError *multierror.Error
			for _, err := range c.Errors {
				multiError = multierror.Append(multiError, err.Err)
			}

			// Log the errors
			log.Println(multiError.Error())

			// Respond with 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": multiError.Error()})
		}
	}
}

func EnsureTrailingSlash() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasSuffix(c.Request.URL.Path, "/") {
			c.Request.URL.Path += "/"
		}
		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		t := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(t)
		log.Print("Latency: ", latency)

		// Access the status we are sending
		status := c.Writer.Status()
		log.Println("Status: ", status)
	}
}
