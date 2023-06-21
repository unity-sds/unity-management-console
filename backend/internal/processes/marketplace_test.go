package processes

import (
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
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				convey.So(req.URL.String(), convey.ShouldEqual, "/applications/test-app/v1.0.0/metadata.json")

				// Send response to be tested
				rw.Write([]byte(`{"name":"test-app","version":"v1.0.0","description":"test"}`))
			}))
			// Close the server when test finishes
			defer server.Close()

			metadata, err := fetchMarketplaceMetadata(name, version)
			convey.Convey("Then there should be no error", func() {
				convey.So(err, convey.ShouldBeNil)
			})
			convey.Convey("And the metadata should be as expected", func() {
				convey.So(metadata.Name, convey.ShouldEqual, "test-app")
				convey.So(metadata.Version, convey.ShouldEqual, "v1.0.0")
			})
		})

		convey.Convey("When the metadata does not exist", func() {
			// Start a local HTTP server to simulate the Github raw content
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				convey.So(req.URL.String(), convey.ShouldEqual, "/applications/test-app/v1.0.0/metadata.json")

				// Send not found response
				http.NotFound(rw, req)
			}))
			// Close the server when test finishes
			defer server.Close()

			metadata, err := fetchMarketplaceMetadata(name, version)
			convey.Convey("Then there should be an error", func() {
				convey.So(err, convey.ShouldNotBeNil)
			})
			convey.Convey("And the metadata should be empty", func() {
				// Assuming your MarketplaceMetadata struct uses pointers for all fields
				convey.So(metadata, convey.ShouldBeZeroValue)
			})
		})
	})
}
