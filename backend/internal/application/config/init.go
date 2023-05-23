package config

import (
	log "github.com/sirupsen/logrus"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever/githubretriever"
	"time"
)

func InitApplication() (f *ffclient.GoFeatureFlag) {
	return initFeatureFlags()
}

func initFeatureFlags() (f *ffclient.GoFeatureFlag) {
	ff, err := ffclient.New(ffclient.Config{
		PollingInterval: 3 * time.Second,
		Retriever: &githubretriever.Retriever{
			RepositorySlug: "unity-sds/unity-management-console",
			Branch:         "main",
			FilePath:       "configs/flag-config.yaml",
			Timeout:        2 * time.Second,
		},
	})
	defer ffclient.Close()

	if err != nil {
		log.Errorf("Error launching feature flag client. %v", err)
	}

	return ff
}
