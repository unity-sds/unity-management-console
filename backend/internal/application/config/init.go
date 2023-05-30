package config

import (
	log "github.com/sirupsen/logrus"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/thomaspoignant/go-feature-flag/retriever/githubretriever"
	"time"
)

type FeatureFlagClient interface {
	Close()
	BoolVariation(s string, u ffuser.User, f bool) (bool, error)
	// Add other methods as required
}

func InitApplication() {
	FFClient = initFeatureFlags()
}

var FFClient FeatureFlagClient

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
