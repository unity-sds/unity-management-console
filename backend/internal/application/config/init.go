package config

import (
	log "github.com/sirupsen/logrus"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/thomaspoignant/go-feature-flag/retriever/githubretriever"
	"time"
)

type SSMParameter struct {
	Name  string `mapstructure:"name"`
	Type  string `mapstructure:"type"`
	Value string `mapstructure:"value"`
}

type MarketplaceItem struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type AppConfig struct {
	MarketplaceBaseUrl   string
	MarketplaceOwner     string
	MarketplaceRepo      string
	AWSRegion            string
	BucketName           string
	Workdir              string
	DefaultSSMParameters []SSMParameter
	BasePath             string
	ConsoleHost          string
	InstallPrefix        string
	Project              string
	Venue                string
	MarketplaceItems     []MarketplaceItem `mapstructure:"MarketplaceItems"` 
	Version				 string
	// GithubToken removed - now using SSM only
}

type FeatureFlagClient interface {
	Close()
	BoolVariation(s string, u ffuser.User, f bool) (bool, error)
	// Add other methods as required
}

func InitApplication() {
	FFClient = initFeatureFlags()
}

var FFClient FeatureFlagClient

// initFeatureFlags initializes a feature flag client by connecting
// to a GitHub repository where the configuration file resides. It
// retrieves configuration settings from a specific file in the
// repository. The repository is 'unity-sds/unity-management-console',
// the branch is 'main', and the file is 'configs/flag-config.yaml'.
//
// The function uses a polling interval of 3 seconds to check for
// updates to the feature flags and a timeout of 2 seconds for
// retrieving the file.
//
// If there is an error during the creation of the feature flag
// client, it will be logged and the function will return a nil
// client.
//
// Example usage:
//
//	featureFlags := initFeatureFlags()
//	if featureFlags == nil {
//	    // Handle the error
//	}
//
// Returns:
//
//	FeatureFlagClient : The initialized feature flag client
func initFeatureFlags() (f FeatureFlagClient) {
	ff, err := ffclient.New(ffclient.Config{
		PollingInterval: 3 * time.Second,
		Retriever: &githubretriever.Retriever{
			RepositorySlug: "unity-sds/unity-management-console",
			Branch:         "main",
			FilePath:       "configs/flag-config.yaml",
			Timeout:        2 * time.Second,
		},
	})

	if err != nil {
		log.Errorf("Error launching feature flag client. %v", err)
	}
	var client FeatureFlagClient = ff
	return client
}
