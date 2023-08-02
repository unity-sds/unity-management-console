package processes

import (
	"fmt"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestFetchMarketplaceMetadata(t *testing.T) {

	convey.Convey("Given a name and version", t, func() {
		name := "test-app"
		version := "v1.0.0"

		convey.Convey("When the metadata exists", func() {
			// Start a local HTTP server to simulate the Github raw content
			var requestURL string
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				requestURL = req.URL.String()

				// Send response to be tested
				rw.Write([]byte(`{"Name":"test-app","Version":"v1.0.0","Description":"test"}`))
			}))
			// Close the server when test finishes
			defer server.Close()
			appConfig := config.AppConfig{MarketplaceBaseUrl: fmt.Sprintf("%s/", server.URL), MarketplaceRepo: "unity-marketplace", MarketplaceOwner: "unity-sds"}

			metadata, err := fetchMarketplaceMetadata(name, version, &appConfig)

			convey.Convey("Then the request URL should be as expected", func() {
				convey.So(requestURL, convey.ShouldEqual, "/unity-sds/unity-marketplace/main/applications/test-app/v1.0.0/metadata.json")
			})

			convey.Convey("Then there should be no error", func() {
				convey.So(err, convey.ShouldBeNil)
			})
			convey.Convey("And the metadata should be as expected", func() {
				convey.So(metadata.Name, convey.ShouldEqual, "test-app")
				convey.So(metadata.Version, convey.ShouldEqual, "v1.0.0")
			})
		})

	})
}

func TestFetchMarketplaceMetadataFailure(t *testing.T) {
	convey.Convey("Given a name and version", t, func() {
		name := "test-app"
		version := "v1.0.0"

		convey.Convey("When the metadata does not exist", func() {
			// Start a local HTTP server to simulate the Github raw content
			var requestURL string
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				requestURL = req.URL.String()

				// Send not found response
				http.NotFound(rw, req)
			}))
			// Close the server when test finishes
			defer server.Close()
			appConfig := config.AppConfig{MarketplaceBaseUrl: fmt.Sprintf("%s/", server.URL), MarketplaceRepo: "unity-marketplace", MarketplaceOwner: "unity-sds"}

			metadata, err := fetchMarketplaceMetadata(name, version, &appConfig)

			convey.Convey("Then the request URL should be as expected", func() {
				convey.So(requestURL, convey.ShouldEqual, "/unity-sds/unity-marketplace/main/applications/test-app/v1.0.0/metadata.json")
			})

			convey.Convey("Then there should be an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
			})
			convey.Convey("And the metadata should be empty", func() {
				// Assuming your MarketplaceMetadata struct uses pointers for all fields
				convey.So(metadata.Name, convey.ShouldBeEmpty)
			})
		})
	})
}
