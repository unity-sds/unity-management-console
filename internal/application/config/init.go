package config

import (
	"time"

	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever/githubretriever"
)

func InitApplication() (f *ffclient.GoFeatureFlag) {
	return initFeatureFlags()
}

func initFeatureFlags() (f *ffclient.GoFeatureFlag) {
	ff, err := ffclient.New(ffclient.Config{
		PollingInterval: 3 * time.Second,
		Retriever: &githubretriever.Retriever{
			RepositorySlug: "unity-sds/unity-control-plane",
			Branch:         "main",
			FilePath:       "configs/flag-config.yaml",
			Timeout:        2 * time.Second,
		},
	})
	defer ffclient.Close()

	if err != nil {

	}
	return ff
}
